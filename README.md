Musescore Scraper

A high-performance web scraper for Musescore sheet music, built with Go and TypeScript. This tool is designed to identify and export score assets for offline study and practice.

🚀 Features

High-Performance Backend: Built with Go, utilizing concurrency for efficient data retrieval.

Modern Web Interface: Includes a frontend component for a visual management of the scraping process.

Container-Ready: Includes configuration files for Docker and automated environments.

Asset Detection: Capable of handling various score formats and image assets (PNG, etc.).

Cloud Compatibility: Pre-configured for deployment on platforms like Railway and Nixpacks.

🛠️ Tech Stack

Core: Go (Golang)

Frontend: TypeScript, JavaScript, CSS

Infrastructure: Docker, Makefile, Nixpacks

Deployment Platform: Railway

⚙️ Project Structure

The repository contains all the necessary configuration for a full-stack Go application:

src/ & dist/: Source code and build outputs.

Dockerfile: Containerization setup.

railway.json & nixpacks.toml: Infrastructure as code for cloud hosting.

Makefile: Task automation for development workflows.

go.mod: Dependency management for the Go ecosystem.

⚠️ Disclaimer

This tool is for educational and personal use only. The scraping of copyrighted material may violate the terms of service of the target website. This project was created to demonstrate web scraping techniques and for personal study purposes. Please support composers, arrangers, and artists by purchasing official sheet music whenever possible. The author does not condone copyright infringement and is not responsible for any misuse of this software.

Developed by FedericoDeniard
