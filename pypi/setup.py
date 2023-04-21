from setuptools import setup, find_packages
import os

with open(os.path.join("..", "README.md"), "r", encoding="utf-8") as fh:
    long_description = fh.read()

setup(
    name='langforge-cli',
    version='0.1.0',
    author="Markus Ecker",
    author_email="markus.ecker@gmail.com",
    description="A Toolkit for Creating and Deploying LangChain Apps",
    long_description=long_description,
    long_description_content_type="text/markdown",
    url="https://github.com/mme/langforge",
    packages=find_packages(),
    classifiers=[
        "Development Status :: 3 - Alpha",
        "Intended Audience :: Developers",
        "License :: OSI Approved :: MIT License",
        "Operating System :: OS Independent",
        "Programming Language :: Python :: 3",
        "Programming Language :: Python :: 3.6",
        "Programming Language :: Python :: 3.7",
        "Programming Language :: Python :: 3.8",
        "Programming Language :: Python :: 3.9",
    ],
    python_requires=">=3.6",
    include_package_data=True,
    package_data={'langforge': ['bin/*']},
    entry_points={
        'console_scripts': [
            'langforge = langforge.main:main',
        ],
    },
)
