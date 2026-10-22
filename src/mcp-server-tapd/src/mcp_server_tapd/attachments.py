"""Display defaults: inline images, attached videos with players, other files attached."""
import os

MAX_BYTES = 150 * 1024 * 1024
IMAGE_EXTENSIONS = {'.png', '.jpg', '.jpeg', '.gif', '.webp', '.bmp'}
VIDEO_EXTENSIONS = {'.mp4', '.mov', '.m4v', '.webm', '.avi', '.mkv'}


def prepare_attachments(data, resolve_path):
    raw = data.pop('attachments', [])
    if not isinstance(raw, list):
        raise ValueError('attachments 必须是数组，元素为 {path, filename}')
    if len(raw) > 20:
        raise ValueError('单次最多上传 20 个附件')
    pending = []
    images = []
    for item in raw:
        if not isinstance(item, dict) or not isinstance(item.get('path'), str):
            raise ValueError('附件必须提供本地 path 或 media:// 引用')
        p = resolve_path(item['path'])
        if '://' in p or not os.path.isfile(p):
            raise ValueError('附件必须是已下载的本地文件')
        name = item.get('filename') or os.path.basename(p)
        display = item.get('display', 'auto')
        if display not in ('auto', 'attachment'):
            raise ValueError('display 必须为 auto 或 attachment')
        if not isinstance(name, str) or not name.strip() or any(c in name for c in '/\\\r\n\x00') or name in {'.', '..'}:
            raise ValueError('filename 必须是原始文件名，不能含目录或换行')
        if os.path.getsize(p) > MAX_BYTES:
            raise ValueError(f'附件 {name} 超过 150MB')
        if display == 'auto' and os.path.splitext(name)[1].lower() in IMAGE_EXTENSIONS:
            if os.path.splitext(p)[1].lower() not in IMAGE_EXTENSIONS:
                raise ValueError(f'图片 {name} 的本地路径没有图片扩展名，请先保存为正确图片类型后内嵌')
            images.append(p)
        else:
            pending.append({'path': p, 'filename': name, 'display': display})
    if raw and data.get('entity_type') == 'tasks':
        raise ValueError('Task 附件尚未验证，请使用需求 stories 或缺陷 bug；未创建/更新工作项')
    if images:
        existing = data.get('image_urls') or []
        if isinstance(existing, str): existing = [existing]
        data['image_urls'] = list(existing) + images
    return pending


def upload_files(client, workspace_id, entity, entry_id, files, render_video=None):
    if entity not in {'story', 'bug'}:
        raise ValueError('附件挂载只支持 story / bug')
    results = []
    previews = []
    for f in files:
        result = {'filename': f['filename'], 'status': 'failed'}
        try:
            if not entry_id:
                raise ValueError('missing entry id')
            if os.path.getsize(f['path']) > MAX_BYTES:
                raise ValueError('attachment exceeds 150MB')
            ret = client.upload_attachment(workspace_id, f['path'], entity, entry_id,
                                           filename=f['filename'])
            attachment_id = ((ret.get('data') or {}).get('Attachment') or {}).get('id')
            if ret.get('status') == 1 and attachment_id:
                result.update(status='uploaded', attachment_id=str(attachment_id))
                if f.get('display', 'auto') == 'auto' and os.path.splitext(f['filename'])[1].lower() in VIDEO_EXTENSIONS:
                    previews.append(result)
            else:
                result['error'] = 'TAPD 未确认上传成功或未返回附件 ID'
        except Exception as exc:
            # No credentials, file contents or server response dumps in error output.
            result['error'] = f'附件未确认上传成功（{type(exc).__name__}）；请先核对票内附件再重试'
        results.append(result)
    if previews:
        append_video_players(client, workspace_id, entity, entry_id, previews, render_video)
    return results


def append_video_players(client, workspace_id, entity, entry_id, results, render_video):
    """Read the existing description before appending; a preview failure never undoes upload."""
    for r in results:
        r['preview_status'] = 'failed'
    try:
        if render_video is None:
            raise ValueError('video renderer unavailable')
        key = 'Bug' if entity == 'bug' else 'Story'
        query = {'workspace_id': workspace_id, 'id': entry_id, 'fields': 'id,description'}
        current = client.get_bug(query) if entity == 'bug' else client.get_stories(query)
        if current.get('status') != 1 or not isinstance(current.get('data'), list):
            raise ValueError('cannot read current description')
        rows = [r[key] for r in current['data'] if isinstance(r, dict) and isinstance(r.get(key), dict) and str(r[key].get('id')) == str(entry_id)]
        if len(rows) != 1 or 'description' not in rows[0]:
            raise ValueError('description missing; refusing to overwrite')
        description = rows[0]['description'] or ''
        if not isinstance(description, str):
            raise ValueError('invalid description')
        blocks = [render_video(workspace_id, r['attachment_id'], entity) for r in results]
        payload = {'workspace_id': workspace_id, 'id': entry_id,
                   'description': '\n'.join([description, *blocks])}
        updated = client.create_or_update_bug(payload) if entity == 'bug' else client.create_or_update_story(payload)
        if updated.get('status') != 1:
            raise ValueError('description update rejected')
        for r in results:
            r['preview_status'] = 'embedded'
    except Exception as exc:
        for r in results:
            r['preview_error'] = f'附件已上传，但正文播放器回填未确认成功（{type(exc).__name__}）；不要重复上传附件'


def finalize_attachments(client, workspace_id, entity, response, files, warnings, render_video=None):
    if not files:
        return response
    key = 'Bug' if entity == 'bug' else 'Story'
    entry_id = ((response.get('data') or {}).get(key) or {}).get('id') if isinstance(response, dict) else None
    if not isinstance(response, dict):
        warnings.append('工作项响应格式不正确，附件未上传')
        return response
    # A failed create/update must not cause a secondary write.
    if response.get('status') != 1:
        entry_id = None
    results = upload_files(client, workspace_id, entity, entry_id, files, render_video=render_video)
    response.setdefault('attachment_results', []).extend(results)
    for r in results:
        if r['status'] != 'uploaded': warnings.append(f"附件 {r['filename']}：{r['error']}")
        elif r.get('preview_status') == 'failed': warnings.append(f"附件 {r['filename']}：{r['preview_error']}")
    return response
