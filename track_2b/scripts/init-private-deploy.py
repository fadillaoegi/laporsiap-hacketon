"""Create private local deployment credentials without printing them."""
import os
import secrets
from pathlib import Path

target = Path(__file__).resolve().parents[1] / 'deploy/.env'
body = 'DEPLOY_DB_PASSWORD=' + secrets.token_hex(32) + '\nDEPLOY_APERTUS_BASE_URL=\nDEPLOY_APERTUS_API_KEY=\n'
descriptor = os.open(target, os.O_WRONLY | os.O_CREAT | os.O_EXCL, 0o600)
with os.fdopen(descriptor, 'w') as output:
    output.write(body)
print('Private deploy/.env created with mode 600. Existing credentials are never overwritten.')
