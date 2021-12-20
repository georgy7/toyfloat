#!/bin/bash
set -e
cd "$(dirname "$0")"

# sudo apt install texlive-latex-extra
# sudo apt install dvipng

latex -output-format=dvi formula.tex
dvipng -bg 'rgb 0.9 0.9 0.9' -o formula.png formula.dvi

latex -output-format=dvi bits.tex
dvipng -bg 'rgb 0.9 0.9 0.9' -o bits.png bits.dvi

# sudo python3 -m pip install matplotlib

go run gen_precision.go
python3 plot.py
