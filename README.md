# LaporSiap — Hack Apertus 2026, Track 2B (Own Project)

LaporSiap helps a person in Indonesia who suspects a part-time-job / top-up
scam turn an untidy chat into reviewed report material for the official
channel (OJK IASC). Apertus drafts a source-linked chronology; the person
checks every entry before anything can be exported.

The project root is [`track_2b/`](track_2b/), following the
[Hack Apertus project template](https://github.com/HackApertus/project-template).

```sh
cd track_2b
make run        # PostgreSQL + LaporSiap in Docker → http://localhost:8787
```

| Path | Contents |
| --- | --- |
| `track_2b/README.md` | Run instructions and project overview |
| `track_2b/technical_report.md` | Architecture, evaluation, limitations |
| `track_2b/Makefile` | `make run` |
| `track_2b/src/` | Go backend (`laporsiap_BE`) and React frontend (`laporsiap_FE`) |
| `track_2b/data/` | Synthetic datasets |
| `track_2b/docs/` | Notes, test scenarios, evaluation fixtures, screenshots |

## License

Code: Apache-2.0 ([LICENSE](LICENSE)). Documentation and other non-code
material: CC-BY-4.0. Synthetic datasets: CDLA-Permissive-2.0. See
[NOTICE](NOTICE) and [LICENSES/](LICENSES/).
