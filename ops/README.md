# Operations Notes

- Build frontend in `apps/web` and copy output to `apps/server/web` before building Go binary.
- Go binary serves embedded static assets and `/api/health`.
- Systemd units in `ops/systemd` are templates and need path/user/display adaptation.
