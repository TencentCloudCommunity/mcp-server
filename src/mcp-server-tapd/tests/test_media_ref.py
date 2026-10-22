"""media:// 引用解析——2026-09-10 那条评论里图片变成灰色占位的根因。"""
import os
import tempfile
import unittest
from unittest.mock import patch

from mcp_server_tapd import server


class ResolveMediaRefTests(unittest.TestCase):
    def setUp(self) -> None:
        self.tmp = tempfile.TemporaryDirectory()
        self.media_dir = self.tmp.name
        os.makedirs(os.path.join(self.media_dir, "inbound"))
        self.img = os.path.join(self.media_dir, "inbound", "3c04ffd3.jpg")
        with open(self.img, "wb") as fh:
            fh.write(b"\xff\xd8\xff\xe0 fake jpeg")
        self._p = patch.object(server, "_OPENCLAW_MEDIA_DIR", self.media_dir)
        self._p.start()

    def tearDown(self) -> None:
        self._p.stop()
        self.tmp.cleanup()

    def test_maps_inbound_ref_to_the_file_on_disk(self) -> None:
        self.assertEqual(
            server._resolve_media_ref("media://inbound/3c04ffd3.jpg"),
            os.path.realpath(self.img),
        )

    def test_leaves_everything_else_alone(self) -> None:
        for url in ("https://x/y.png", "/tfl/pictures/a.png", "/tmp/a.png", "", None):
            self.assertEqual(server._resolve_media_ref(url), url)

    def test_missing_file_raises_instead_of_embedding_a_dead_link(self) -> None:
        with self.assertRaises(ValueError) as ctx:
            server._resolve_media_ref("media://inbound/nope.jpg")
        self.assertIn("找不到", str(ctx.exception))

    def test_rejects_path_traversal(self) -> None:
        for ref in ("media://inbound/../../etc/passwd", "media://../x.png", "media://./a.png"):
            with self.assertRaises(ValueError):
                server._resolve_media_ref(ref)


class ExtractMediaItemsResolvesRefsTests(unittest.TestCase):
    def setUp(self) -> None:
        self.tmp = tempfile.TemporaryDirectory()
        os.makedirs(os.path.join(self.tmp.name, "inbound"))
        self.img = os.path.join(self.tmp.name, "inbound", "a.jpg")
        open(self.img, "wb").write(b"x")
        self._p = patch.object(server, "_OPENCLAW_MEDIA_DIR", self.tmp.name)
        self._p.start()

    def tearDown(self) -> None:
        self._p.stop()
        self.tmp.cleanup()

    # 复现真实调用：image_url 传的就是 media://inbound/<id>
    def test_image_url_media_ref_becomes_a_local_path_so_it_gets_uploaded(self) -> None:
        data = {"image_url": "media://inbound/a.jpg"}
        items = server._extract_media_items(data)
        self.assertEqual(items, [{"type": "image", "url": os.path.realpath(self.img), "alt": ""}])
        self.assertTrue(server._is_local_image_path(items[0]["url"]))

    def test_image_urls_and_media_list_are_resolved_too(self) -> None:
        data = {
            "image_urls": ["media://inbound/a.jpg", "https://cdn/x.png"],
            "media": [{"type": "image", "url": "media://inbound/a.jpg"}],
        }
        urls = [i["url"] for i in server._extract_media_items(data)]
        self.assertEqual(urls.count(os.path.realpath(self.img)), 2)
        self.assertIn("https://cdn/x.png", urls)


if __name__ == "__main__":
    unittest.main()
