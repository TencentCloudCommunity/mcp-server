"""视频内嵌的纯逻辑：哪些算本地视频、怎么分流、HTML 拼成什么样。

不碰网络：上传和建票都要真凭据，这里只锁住"接口决定的形状"——
真正会悄悄坏掉的是这部分，而不是 requests.post。
"""

from __future__ import annotations

import os
import sys
import tempfile
import unittest
from pathlib import Path
from unittest.mock import patch

sys.path.insert(0, str(Path(__file__).resolve().parents[1] / "src"))

from mcp_server_tapd.server import (  # noqa: E402
    _attach_local_videos,
    _human_size,
    _is_local_video_path,
    _render_tapd_video_html,
    _split_local_videos,
)


class IsLocalVideoPathTests(unittest.TestCase):
    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.root = Path(self._tmp.name)
        self.addCleanup(self._tmp.cleanup)

    def _file(self, name: str) -> str:
        p = self.root / name
        p.write_bytes(b"\x00")
        return str(p)

    def test_accepts_a_local_mp4(self) -> None:
        self.assertTrue(_is_local_video_path(self._file("rec.mp4")))

    def test_accepts_other_video_extensions(self) -> None:
        for name in ("a.mov", "a.m4v", "a.webm", "a.avi", "a.mkv"):
            with self.subTest(name=name):
                self.assertTrue(_is_local_video_path(self._file(name)))

    def test_is_case_insensitive(self) -> None:
        self.assertTrue(_is_local_video_path(self._file("REC.MP4")))

    def test_rejects_a_url(self) -> None:
        # 已经是直链的原样放行，不该被当成待上传
        self.assertFalse(_is_local_video_path("https://example.com/a.mp4"))
        self.assertFalse(_is_local_video_path("//cdn.example.com/a.mp4"))

    def test_rejects_a_nonexistent_path(self) -> None:
        self.assertFalse(_is_local_video_path(str(self.root / "missing.mp4")))

    def test_rejects_a_non_video_file(self) -> None:
        self.assertFalse(_is_local_video_path(self._file("shot.png")))

    def test_rejects_empty(self) -> None:
        self.assertFalse(_is_local_video_path(""))


class SplitLocalVideosTests(unittest.TestCase):
    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.local = str(Path(self._tmp.name) / "rec.mp4")
        Path(self.local).write_bytes(b"\x00")
        self.addCleanup(self._tmp.cleanup)

    def test_pulls_local_videos_out_and_leaves_the_rest(self) -> None:
        items = [
            {"type": "image", "url": "/tmp/a.png"},
            {"type": "video", "url": self.local},
            {"type": "video", "url": "https://example.com/remote.mp4"},
        ]

        remaining, videos = _split_local_videos(items)

        self.assertEqual(videos, [self.local])
        # 远程视频留在原地：它不需要上传，直接内嵌就行
        self.assertEqual(
            [i["url"] for i in remaining], ["/tmp/a.png", "https://example.com/remote.mp4"]
        )

    def test_leaves_everything_alone_when_there_are_no_local_videos(self) -> None:
        items = [{"type": "image", "url": "/tmp/a.png"}]
        remaining, videos = _split_local_videos(items)
        self.assertEqual(videos, [])
        self.assertEqual(remaining, items)

    def test_handles_an_empty_list(self) -> None:
        self.assertEqual(_split_local_videos([]), ([], []))


class RenderTapdVideoHtmlTests(unittest.TestCase):
    """模板是从网页端建的真票里逆出来的，改动会让票里播不了，所以锁死。"""

    def _html(self, attach_type: str = "bug") -> str:
        with patch.dict(os.environ, {"TAPD_BASE_URL": "https://www.tapd.cn"}, clear=False):
            return _render_tapd_video_html(64984633, "1164984633001000202", attach_type)

    def test_matches_the_shape_tapd_itself_writes(self) -> None:
        got = self._html()
        src = (
            "https://www.tapd.cn/64984633/attachments/preview_attachments/"
            "1164984633001000202/bug"
        )
        poster = (
            "https://www.tapd.cn/64984633/attachments/preview_video_poster/1164984633001000202"
        )
        self.assertIn(f'src="{src}"', got)
        self.assertIn(f'poster="{poster}"', got)
        self.assertIn(f'<source src="{src}">', got)
        self.assertIn('controls="controls"', got)
        # 打开票时不该去拉视频
        self.assertIn('preload="none"', got)

    def test_url_carries_the_attachment_type(self) -> None:
        # 路径末段就是 type，story 走另一个挂载点
        got = self._html("story")
        self.assertIn("/1164984633001000202/story", got)

    def test_src_is_stable_not_a_signed_download_url(self) -> None:
        # 用 download_url 会过期，票过几天就播不了
        got = self._html()
        self.assertNotIn("?sign", got)
        self.assertNotIn("file.tapd", got)


if __name__ == "__main__":  # pragma: no cover
    unittest.main()


class SilentDropGuardTests(unittest.TestCase):
    """不支持内嵌视频的接口（评论）必须报错，不能悄悄丢。"""

    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.local = str(Path(self._tmp.name) / "rec.mp4")
        Path(self.local).write_bytes(b"\x00")
        self.addCleanup(self._tmp.cleanup)

    def test_raises_when_the_caller_cannot_handle_local_video(self) -> None:
        from mcp_server_tapd.server import _render_rich_description

        data = {"workspace_id": 1, "description": "见录屏", "video_url": self.local}
        with self.assertRaisesRegex(ValueError, "不支持内嵌本地视频"):
            _render_rich_description(data)

    def test_returns_the_paths_when_the_caller_can(self) -> None:
        from mcp_server_tapd.server import _render_rich_description

        data = {"workspace_id": 1, "description": "见录屏", "video_url": self.local}
        pending = _render_rich_description(data, supports_local_video=True)
        self.assertEqual(pending, [self.local])
        # 本地视频这一步不该进描述——它要等票建出来才有 attachment id
        self.assertNotIn(self.local, data["description"])

    def test_a_remote_video_still_renders_inline_immediately(self) -> None:
        from mcp_server_tapd.server import _render_rich_description

        data = {"workspace_id": 1, "description": "见录屏",
                "video_url": "https://example.com/a.mp4"}
        pending = _render_rich_description(data)
        self.assertEqual(pending, [])
        self.assertIn("<video", data["description"])
        self.assertIn("https://example.com/a.mp4", data["description"])


class OversizeIsAWarningNotABlockerTests(unittest.TestCase):
    """票必须建出来。超限只是少一个附件，不该把整张票挡下来。

    尤其视频是建票**之后**才传的：抛异常的结果是"票建了、调用方收到失败"，
    调用方多半会重试，于是多一张重复的票。
    """

    def setUp(self) -> None:
        self._tmp = tempfile.TemporaryDirectory()
        self.big = Path(self._tmp.name) / "huge.mp4"
        self.big.write_bytes(b"\x00" * 1024)
        self.addCleanup(self._tmp.cleanup)

    def test_oversize_video_warns_and_is_skipped(self) -> None:
        warnings: list[str] = []
        with patch("mcp_server_tapd.server._MAX_ATTACHMENT_BYTES", 512):
            html_out = _attach_local_videos(1, "bug", 99, [str(self.big)], warnings)

        self.assertEqual(html_out, "")  # 描述里没有这个视频
        self.assertEqual(len(warnings), 1)
        self.assertIn("huge.mp4", warnings[0])
        self.assertIn("未附到票上", warnings[0])

    def test_upload_failure_warns_rather_than_raising(self) -> None:
        warnings: list[str] = []
        with patch("mcp_server_tapd.server.client") as fake:
            fake.upload_attachment.side_effect = RuntimeError("boom")
            html_out = _attach_local_videos(1, "bug", 99, [str(self.big)], warnings)

        self.assertEqual(html_out, "")
        self.assertIn("boom", warnings[0])

    def test_api_says_failed_warns_rather_than_raising(self) -> None:
        warnings: list[str] = []
        with patch("mcp_server_tapd.server.client") as fake:
            fake.upload_attachment.return_value = {"status": 0, "info": "quota exceeded"}
            html_out = _attach_local_videos(1, "bug", 99, [str(self.big)], warnings)

        self.assertEqual(html_out, "")
        self.assertIn("quota exceeded", warnings[0])

    def test_one_bad_video_does_not_lose_the_good_one(self) -> None:
        good = Path(self._tmp.name) / "ok.mp4"
        good.write_bytes(b"\x00")
        warnings: list[str] = []
        with patch("mcp_server_tapd.server.client") as fake, \
                patch("mcp_server_tapd.server._MAX_ATTACHMENT_BYTES", 512):
            fake.upload_attachment.return_value = {
                "status": 1, "data": {"Attachment": {"id": "777"}}
            }
            html_out = _attach_local_videos(1, "bug", 99, [str(self.big), str(good)], warnings)

        self.assertIn("777", html_out)          # 小的那个贴上去了
        self.assertEqual(len(warnings), 1)      # 大的那个只记了一条警告

    def test_oversize_image_warns_when_a_collector_is_given(self) -> None:
        from mcp_server_tapd.server import _upload_local_images

        img = Path(self._tmp.name) / "shot.png"
        img.write_bytes(b"\x00" * 1024)
        warnings: list[str] = []
        with patch("mcp_server_tapd.server._MAX_UPLOAD_BYTES", 512):
            out = _upload_local_images(
                [{"type": "image", "url": str(img)}], 1, warnings
            )

        self.assertEqual(out, [])               # 这张图不进描述
        self.assertIn("shot.png", warnings[0])


class HumanSizeTests(unittest.TestCase):
    def test_reports_megabytes(self) -> None:
        self.assertEqual(_human_size(8 * 1024 * 1024), "8.0MB")

    def test_falls_back_to_kilobytes(self) -> None:
        self.assertEqual(_human_size(200 * 1024), "200KB")


class UploadTypeTests(unittest.TestCase):
    """上传接口只认实体名本身。

    2026-09-09 实测：story ✅；story_description / story_attachment /
    story_description_attachment / stories / workitem 全部 422 "type is invalid"。
    第一版从 GET /attachments 的返回里抄了 bug_description，每次 422，
    而 422 当时又被吞掉，表现成"没报错但也没传上"。
    """

    def test_uses_the_bare_entity_name(self) -> None:
        from mcp_server_tapd.server import _VIDEO_ATTACH_TYPE

        self.assertEqual(_VIDEO_ATTACH_TYPE["bug"], "bug")
        self.assertEqual(_VIDEO_ATTACH_TYPE["story"], "story")

    def test_never_the_description_variants(self) -> None:
        from mcp_server_tapd.server import _VIDEO_ATTACH_TYPE

        for value in _VIDEO_ATTACH_TYPE.values():
            self.assertNotIn("description", value)
            self.assertNotIn("attachment", value)

    def test_upload_is_called_with_that_type(self) -> None:
        import tempfile as tf
        from unittest.mock import patch as _patch

        from mcp_server_tapd.server import _attach_local_videos

        with tf.TemporaryDirectory() as d:
            v = Path(d) / "rec.mp4"
            v.write_bytes(b"\x00")
            with _patch("mcp_server_tapd.server.client") as fake:
                fake.upload_attachment.return_value = {
                    "status": 1, "data": {"Attachment": {"id": "1"}}
                }
                _attach_local_videos(64984633, "story", "S1", [str(v)], [])
            self.assertEqual(fake.upload_attachment.call_args.args[2], "story")
