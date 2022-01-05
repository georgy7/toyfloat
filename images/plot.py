import matplotlib.pyplot as plt
import csv
import numpy as np

def get_data(data_fn, max_value):
    values = []
    precision = []

    tsv_file = open(data_fn)
    read_tsv = csv.reader(tsv_file, delimiter="\t")

    for row in read_tsv:
        v = float(row[0])
        p = float(row[1])

        if (v < max_value) and (p > 1e-7):
            values.append(v)
            precision.append(p)

    return values, precision


def get_comparison_data(data_fn, minPowerX10 = -22, maxPowerX10 = 21):
    values, precision = get_data(data_fn, 20000)
    assert len(values) == len(precision)

    rv, pv = [], []

    for power10 in range(minPowerX10, maxPowerX10, 2):
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


def make_image(data_fn, result_fn, max_value = 20000):
    values, precision = get_data(data_fn, max_value)

    figure, ax = new_figure()
    figure.set_size_inches(5, 5)
    plt.subplots_adjust(left=0.162, bottom=0.129, right=0.954, top=0.954)

    plt.plot(values, precision, '.')
    plt.savefig(result_fn)


def make_comparison(result_fn):
    figure, ax = new_figure()
    figure.set_size_inches(5.6, 5.5)
    plt.subplots_adjust(left=0.162, bottom=0.129, right=0.954, top=0.954)

    values, precision = get_comparison_data('precision8x3.tsv', minPowerX10=-20, maxPowerX10=8)
    plt.plot(values, precision, '-', label='8x3')

    values, precision = get_comparison_data('precision12.tsv')
    plt.plot(values, precision, '-', label='12')

    values, precision = get_comparison_data('precision13.tsv')
    plt.plot(values, precision, '-', label='12u / 13')

    values, precision = get_comparison_data('precision14.tsv')
    plt.plot(values, precision, '-', label='14')

    values, precision = get_comparison_data('precision15x3.tsv', minPowerX10=-20, maxPowerX10=8)
    plt.plot(values, precision, '-', label='15x3')

    values, precision = get_comparison_data('precision15x2.tsv', minPowerX10=-9, maxPowerX10=3)
    plt.plot(values, precision, '-', color='tab:pink', label='15x2')

#     values, precision = get_comparison_data('precision16u.tsv')
#     plt.plot(values, precision, '-', label='16u')

    values, precision = get_comparison_data('precision16x3u.tsv', minPowerX10=-20, maxPowerX10=8)
    plt.plot(values, precision, '-', color='tab:gray', label='16x3u')

    ax.legend()

    plt.savefig(result_fn)


if __name__ == "__main__":
    make_image('precision12.tsv', 'precision12.png')
    make_image('precision12u.tsv', 'precision12u.png')
    make_image('precision13.tsv', 'precision13.png')
    make_image('precision14.tsv', 'precision14.png')
    make_image('precision15x3.tsv', 'precision15x3.png', 2.0)
    make_image('precision15x2.tsv', 'precision15x2.png', 2.0)
    make_image('precision8x3.tsv', 'precision8x3.png')
    make_image('precision4x3u.tsv', 'precision4x3u.png')
    make_image('precision16.tsv', 'precision16.png')
    make_image('precision16u.tsv', 'precision16u.png')
    make_image('precision16x3u.tsv', 'precision16x3u.png')

    make_comparison('comparison.png')
