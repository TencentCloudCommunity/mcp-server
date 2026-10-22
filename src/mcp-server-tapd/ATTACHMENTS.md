# Ticket attachments

`create_story_or_task`, `update_story_or_task`, `create_bug`, and `update_bug`
accept an `options.attachments` array:

```json
{
  "attachments": [
    {"path": "media://inbound/recording.mp4", "filename": "Recording.mp4"},
    {"path": "/local/archive.zip", "filename": "Documents.zip"},
    {"path": "/local/video.mp4", "filename": "Download-only.mp4", "display": "attachment"}
  ]
}
```

`display` defaults to `auto`. Classification uses the original filename extension,
case-insensitively. PNG/JPG/JPEG/GIF/WebP/BMP images go inline via the image API.
MP4/MOV/M4V/WebM/AVI/MKV videos upload as attachments, then add players to the
current ticket description. Other extensions, including unknown ones, upload only
as attachments. `display: "attachment"` disables inline rendering for any file.
The original `video_url` / `video_urls` inputs retain automatic players and use
the same finalization path. No archive extraction or email parsing is performed.

`attach_files_to_ticket` adds files to an existing `story` or `bug`, including
automatic video players and attachment-only images. To inline images on an
existing ticket, use `update_story_or_task` or `update_bug`; the standalone file
tool rejects auto-inline images before uploading anything.

Every ordinary attachment result contains `status` and, when uploaded,
`attachment_id`. Automatic video results additionally report `preview_status`
(`embedded` or `failed`). Failed previews retain successfully uploaded attachment
IDs: do not create another ticket or reupload those files. Description read/update
failures do not turn a successful file upload into an upload failure. Players are
TAPD HTML embeds; actual playback still depends on TAPD and browser codec support.

Only stories and bugs support this path. Task attachment requests are rejected.
Limits: 20 files per request, 150 MiB per attachment, 5 MiB for inline images.
The WeChat download limit is separate. Inputs must be local files or valid
`media://` references; filenames must not contain path separators or newlines.

Tests: `PYTHONPATH=src python -m unittest discover -s tests -v`.
