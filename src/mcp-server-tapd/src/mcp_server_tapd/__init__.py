import asyncio

from . import server
from . import tapd

__version__ = "2.1.4"


def main():
    """Main entry point for the package."""
    asyncio.run(server.main())


__all__ = ["main", "server", "tapd", "__version__"]
