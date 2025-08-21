from . import server
def main():
    """MCP Fetch Server - HTTP fetching functionality for MCP"""
    import asyncio
    asyncio.run(server.serve())


if __name__ == "__main__":
    main()
    

__all__ = ["main", "server"]