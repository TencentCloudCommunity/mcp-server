import json
import tempfile
import unittest
from pathlib import Path
from unittest.mock import Mock, patch
from mcp_server_tapd import server
from mcp_server_tapd.attachments import prepare_attachments, upload_files


class DisplayPolicyTests(unittest.TestCase):
    def setUp(self):
        self.temp = tempfile.TemporaryDirectory()
        self.addCleanup(self.temp.cleanup)
        self.file = Path(self.temp.name) / 'uuid.bin'
        self.file.write_bytes(b'video bytes')
        self.client = Mock()
        self.client.upload_attachment.return_value = {'status': 1, 'data': {'Attachment': {'id': '12'}}}
        self.client.get_stories.return_value = {'status': 1, 'data': [{'Story': {'id': '1', 'description': '<p>Original</p>'}}]}
        self.client.get_bug.return_value = {'status': 1, 'data': [{'Bug': {'id': '1', 'description': '<p>Original</p>'}}]}
        self.client.create_or_update_story.return_value = {'status': 1, 'data': {'Story': {'id': '1'}}}
        self.client.create_or_update_bug.return_value = {'status': 1, 'data': {'Bug': {'id': '1'}}}

    def prepare(self, name='Recording.MP4', **extra):
        return prepare_attachments({'attachments': [{'path': str(self.file), 'filename': name, **extra}]}, lambda p: p)

    def test_auto_video_preserves_original_filename_and_description(self):
        result = upload_files(self.client, 1, 'story', '1', self.prepare(), server._render_tapd_video_html)
        self.assertEqual(result[0]['preview_status'], 'embedded')
        self.assertEqual(self.client.upload_attachment.call_args.kwargs['filename'], 'Recording.MP4')
        description = self.client.create_or_update_story.call_args.args[0]['description']
        self.assertTrue(description.startswith('<p>Original</p>'))
        self.assertIn('<video ', description)

    def test_attachment_override_and_unknown_formats_never_add_player(self):
        for name, options in [('movie.mp4', {'display':'attachment'}), ('image.png', {'display':'attachment'}), ('custom.xyz', {}), ('mail.eml', {})]:
            with self.subTest(name=name):
                client = Mock()
                client.upload_attachment.return_value = {'status':1, 'data':{'Attachment':{'id':'12'}}}
                result = upload_files(client,1,'story','1',self.prepare(name, **options),server._render_tapd_video_html)
                self.assertEqual(result[0]['status'],'uploaded')
                self.assertNotIn('preview_status',result[0])
                client.get_stories.assert_not_called()
                client.create_or_update_story.assert_not_called()

    def test_preview_failure_keeps_uploaded_id_and_does_not_reupload(self):
        self.client.create_or_update_story.return_value = {'status':0}
        result = upload_files(self.client,1,'story','1',self.prepare(),server._render_tapd_video_html)
        self.assertEqual(result[0]['status'],'uploaded')
        self.assertEqual(result[0]['attachment_id'],'12')
        self.assertEqual(result[0]['preview_status'],'failed')
        self.client.upload_attachment.assert_called_once()

    def test_missing_description_prevents_overwriting(self):
        self.client.get_stories.return_value = {'status':1,'data':[{'Story':{'id':'1'}}]}
        result = upload_files(self.client,1,'story','1',self.prepare(),server._render_tapd_video_html)
        self.assertEqual(result[0]['preview_status'],'failed')
        self.client.create_or_update_story.assert_not_called()

    def test_failed_upload_never_attempts_player(self):
        self.client.upload_attachment.return_value = {'status':0}
        result = upload_files(self.client,1,'story','1',self.prepare(),server._render_tapd_video_html)
        self.assertEqual(result[0]['status'],'failed')
        self.client.get_stories.assert_not_called()

    def test_invalid_display_rejected_before_upload(self):
        with self.assertRaises(ValueError): self.prepare(display='inline')

    def test_legacy_and_new_video_paths_share_results(self):
        video = self.file.with_suffix('.mp4'); video.write_bytes(b'video')
        with patch.object(server,'client',self.client):
            ret = server.create_story_or_task(1,'test',{'video_url':str(video), 'attachments':[{'path':str(self.file),'filename':'another.mp4'}]})
        results = json.loads(ret['data'])['attachment_results']
        self.assertEqual(len(results),2)
        self.assertTrue(all(r['preview_status']=='embedded' for r in results))
        self.assertEqual(self.client.upload_attachment.call_count,2)

    def test_existing_bug_tool_uses_default_video_policy(self):
        with patch.object(server,'client',self.client):
            ret=server.attach_files_to_ticket(1,'1','bug',[{'path':str(self.file),'filename':'recording.mp4'}])
        self.assertTrue(ret['ok'])
        self.assertEqual(ret['attachment_results'][0]['preview_status'],'embedded')
        self.assertIn('<video ',self.client.create_or_update_bug.call_args.args[0]['description'])
