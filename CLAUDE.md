# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Repository Structure

This is a monorepo with the following top-level directories:

- `main/` - Main application module
- `db-migrater/` - Database migration tooling
- `docs/` - Project documentation

## Development Workflow

### Branching

- Always create a dedicated branch before starting any work (`feat-<name>`, `fix-<name>`, etc.)
- **Direct commits to `main` are strictly prohibited**

### Before Committing

- Verify the change does not break other parts of the system
- Confirm the code runs without errors before committing
