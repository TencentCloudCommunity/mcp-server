import json
import os
import tempfile
import unittest
from pathlib import Path
from unittest.mock import Mock, patch
from mcp_server_tapd import server
from mcp_server_tapd.attachments import prepare_attachments, upload_files
from mcp_server_tapd.tapd import TAPDClient

class AttachmentsTests(unittest.TestCase):
    def setUp(self):
        self.tmp = tempfile.TemporaryDirectory(dir=os.environ.get('TMPDIR'))
        self.root = Path(self.tmp.name)
        for name in ['uuid.bin', 'image.png', 'archive.zip']:
            (self.root/name).write_bytes(b'test file')

    def tearDown(self): self.tmp.cleanup()

    def files(self): return [{'path':str(self.root/'uuid.bin'),'filename':'邮件线程.eml'}, {'path':str(self.root/'archive.zip'),'filename':'材料.zip'}]

    def test_images_route_inline_and_documents_keep_original_name(self):
        data={'attachments':self.files()+[{'path':str(self.root/'image.png'),'filename':'截图.png'}]}
        files=prepare_attachments(data, lambda x:x)
        self.assertNotIn('attachments',data)
        self.assertEqual(data['image_urls'],[str(self.root/'image.png')])
        self.assertEqual(files[0]['filename'],'邮件线程.eml')

    def test_create_then_upload_and_partial_failure_is_reported(self):
        client=Mock()
        def create(data):
            self.assertNotIn('attachments',data)
            self.assertEqual(client.upload_attachment.call_count,0)
            return {'status':1,'data':{'Story':{'id':'123'}}}
        client.create_or_update_story.side_effect=create
        client.upload_attachment.side_effect=[{'status':1,'data':{'Attachment':{'id':'99'}}}, RuntimeError('private-token')]
        with patch.object(server,'client',client):
            ret=server.create_story_or_task(1,'test',{'entity_type':'stories','attachments':self.files()})
        details=json.loads(ret['data'])['attachment_results']
        self.assertEqual([r['status'] for r in details],['uploaded','failed'])
        self.assertIn('warnings',ret)
        self.assertNotIn('private-token',json.dumps(ret))
        self.assertEqual(client.upload_attachment.call_args_list[0].kwargs['filename'],'邮件线程.eml')

    def test_invalid_input_and_tasks_rejected_before_writes(self):
        with patch.object(server,'client') as client:
            for options in [{'entity_type':'tasks','attachments':self.files()}, {'attachments':[{'path':str(self.root/'uuid.bin'),'filename':'../x'}]}, {'attachments':[{'path':'/missing'}]}]:
                with self.assertRaises(ValueError):server.create_story_or_task(1,'test',options)
            client.create_or_update_story.assert_not_called()

    def test_existing_ticket_images_rejected_without_upload(self):
        with patch.object(server,'client') as client:
            with self.assertRaises(ValueError):server.attach_files_to_ticket(1,'123','bug',[{'path':str(self.root/'image.png')}])
            client.upload_attachment.assert_not_called()

    def test_bug_creation_and_updates(self):
        for fn,args,key in [(server.create_bug,(1,'bug'),'Bug'),(server.update_bug,(1,),'Bug'),(server.update_story_or_task,(1,),'Story')]:
            with self.subTest(function=fn.__name__),patch.object(server,'client') as client:
                response={'status':1,'data':{key:{'id':'456'}}}
                client.create_or_update_bug.return_value=response
                client.create_or_update_story.return_value=response
                client.upload_attachment.return_value={'status':1,'data':{'Attachment':{'id':'789'}}}
                ret=fn(*args,options={'id':'456','attachments':self.files()})
                if isinstance(ret,str):ret=json.loads(ret)
                self.assertEqual(len(json.loads(ret['data'])['attachment_results']),2)
                self.assertEqual(client.upload_attachment.call_args.args[3],'456')

    def test_multipart_filename_and_file_bytes(self):
        client=TAPDClient();client.base_url='https://example.invalid';client.headers={'Content-Type':'application/json'}
        def post(url,**kw):
            self.assertNotIn('Content-Type',kw['headers'])
            name,fh=kw['files']['file'];self.assertEqual(name,'邮件.eml');self.assertEqual(fh.read(),b'test file')
            self.assertEqual(kw['data']['type'],'story')
            res=Mock();res.json.return_value={'status':1};return res
        with patch('mcp_server_tapd.tapd.requests.post',side_effect=post):
            client.upload_attachment(1,str(self.root/'uuid.bin'),'story','123',filename='邮件.eml')

    def test_no_upload_without_successful_ticket_and_size_limit(self):
        from mcp_server_tapd.attachments import finalize_attachments
        client=Mock();warnings=[]
        finalize_attachments(client,1,'story',{'status':0,'data':{}},self.files(),warnings)
        client.upload_attachment.assert_not_called();self.assertEqual(len(warnings),2)
        with patch('mcp_server_tapd.attachments.os.path.getsize',return_value=151*1024*1024):
            with self.assertRaises(ValueError):prepare_attachments({'attachments':self.files()},lambda x:x)

if __name__=='__main__':unittest.main()
