#!/usr/bin/env bash
set -euo pipefail

# Demo mínima de cores ANSI puras para testar preservação de cores
# dentro do viewport preformatado do Shantilly.

# RED
printf '\033[31mRED\033[0m normal\n'

# YELLOW
printf '\033[33mYELLOW\033[0m normal\n'

# GREEN
printf '\033[32mGREEN\033[0m normal\n'

printf '\n[ansi-colors-demo] Fim da saída ANSI.\n'
