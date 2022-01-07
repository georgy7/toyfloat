#!/bin/bash
set -e
cd "$(dirname "$0")"

# sudo apt install texlive-latex-extra
# sudo apt install dvipng

latex -output-format=dvi formula.tex
dvipng -D 130 -bg 'rgb 0.9 0.9 0.9' -o formula.png formula.dvi

# sudo python3 -m pip install matplotlib

go run gen_precision.go
python3 plot.py
