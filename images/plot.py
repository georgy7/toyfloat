import matplotlib.pyplot as plt
import csv
import numpy as np

def get_data(data_fn):
    values = []
    precision = []

    tsv_file = open(data_fn)
    read_tsv = csv.reader(tsv_file, delimiter="\t")

    for row in read_tsv:
        v = float(row[0])
        p = float(row[1])

        if v < 20000:
            values.append(v)
            precision.append(p)

    return values, precision


def get_comparison_data(data_fn):
    values, precision = get_data(data_fn)
    assert len(values) == len(precision)

    rv, pv = [], []

    for power10 in range(-25, 20, 2):
        power = power10 / 10
        min_value = pow(10, power-0.2)
        max_value = pow(10, power)

        max_error = None

        for i in range(len(values)):
            v = values[i]
            e = precision[i]
            if (min_value < v) and (v <= max_value):
                if (max_error is None) or (max_error < e):
                    max_error = e

        if max_error != None:
            rv.append(max_value)
            pv.append(max_error)

    # I want to smooth the graph in such a way that
    # the values on the smoothed graph do not decrease.

    x = np.log10(rv)
    y = np.log10(pv)
    coefficients = np.polyfit(x, y, 1)

    b_fix = 0
    for i in range(len(y)):
        calculated = coefficients[0] * x[i] + coefficients[1]
        if calculated + b_fix < y[i]:
            b_fix = y[i] - calculated

    coefficients[1] += b_fix

    approximated = []
    for i in range(len(rv)):
        calculated = coefficients[0] * x[i] + coefficients[1]
        approximated.append(pow(10, calculated))

    return rv, approximated


def new_figure():
    figure = plt.figure()
    axes = plt.axes()

    figure.patch.set_facecolor('#e6e6e6')
    axes.set_facecolor('#eeeeee')

    plt.xlabel('Floating Point Values')
    plt.ylabel('Floating Point Precision')

    plt.xscale('log')
    plt.yscale('log')

    plt.grid(True, color='#c7c7c7')
    return figure, axes


def make_image(data_fn, result_fn):
    values, precision = get_data(data_fn)

    figure, ax = new_figure()
    figure.set_size_inches(5, 5)
    plt.subplots_adjust(left=0.162, bottom=0.129, right=0.954, top=0.954)

    plt.plot(values, precision, '.')
    plt.savefig(result_fn)


def make_comparison(result_fn):
    figure, ax = new_figure()
    figure.set_size_inches(5.6, 5)
    plt.subplots_adjust(left=0.162, bottom=0.129, right=0.954, top=0.954)

    values, precision = get_comparison_data('precision.tsv')
    plt.plot(values, precision, 'o-', label='default')

    values, precision = get_comparison_data('precision13.tsv')
    plt.plot(values, precision, 'o-', label='unsigned / 13-bit')

    values, precision = get_comparison_data('precision14.tsv')
    plt.plot(values, precision, 'o-', label='14-bit')

    ax.legend()

    plt.savefig(result_fn)


if __name__ == "__main__":
    make_image('precision.tsv', 'precision.png')
    make_image('precision_unsigned.tsv', 'precision_unsigned.png')
    make_image('precision13.tsv', 'precision13.png')
    make_image('precision14.tsv', 'precision14.png')
    make_comparison('comparison.png')
