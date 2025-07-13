#!/usr/bin/env python3

from setuptools import setup, find_packages
import pathlib

here = pathlib.Path(__file__).parent.resolve()

# Get the long description from the README file
long_description = (here / "README.md").read_text(encoding="utf-8")

setup(
    name="urlshortener-sdk",
    version="1.0.0",
    description="Official Python SDK for URL Shortener API",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/yourusername/url-shortener",
    author="URL Shortener Team",
    author_email="support@urlshortener.com",
    classifiers=[
        "Development Status :: 5 - Production/Stable",
        "Intended Audience :: Developers",
        "Topic :: Software Development :: Libraries :: Python Modules",
        "Topic :: Internet :: WWW/HTTP",
        "License :: OSI Approved :: MIT License",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
        "Programming Language :: Python :: 3.10",
        "Programming Language :: Python :: 3.11",
        "Programming Language :: Python :: 3.12",
        "Programming Language :: Python :: 3 :: Only",
    ],
    keywords="url-shortener, api, sdk, python, client",
    packages=find_packages(exclude=["tests*"]),
    python_requires=">=3.8",
    install_requires=[
        "requests>=2.25.0",
        "pydantic>=2.0.0",
        "typing-extensions>=4.0.0; python_version<'3.10'",
    ],
    extras_require={
        "dev": [
            "pytest>=7.0.0",
            "pytest-asyncio>=0.21.0",
            "pytest-cov>=4.0.0",
            "pytest-mock>=3.10.0",
            "black>=23.0.0",
            "isort>=5.12.0",
            "flake8>=6.0.0",
            "mypy>=1.0.0",
            "pre-commit>=3.0.0",
        ],
        "async": [
            "aiohttp>=3.8.0",
            "asyncio-throttle>=1.0.0",
        ],
    },
    project_urls={
        "Bug Reports": "https://github.com/yourusername/url-shortener/issues",
        "Source": "https://github.com/yourusername/url-shortener/tree/main/sdk/python",
        "Documentation": "https://docs.urlshortener.com/sdk/python",
    },
    entry_points={
        "console_scripts": [
            "urlshortener=urlshortener.cli:main",
        ],
    },
)