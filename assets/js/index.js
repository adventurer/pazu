$(function() {
    MorrisArea();

    // progress bars
    $('.progress .progress-bar').progressbar({
            display_text: 'none'
    });

    $('.sparkbar').sparkline('html', { type: 'bar' });

    $('.sparkline-pie').sparkline('html', {
        type: 'pie',
        offset: 90,
        width: '100px',
        height: '100px',
        sliceColors: ['#29bd73', '#182973', '#ffcd55']
    })


    // top products
    var dataStackedBar = {
            labels: ['Q1','Q2','Q3','Q4','Q5'],
            series: [
                [2350,3205,4520,2351,5632],
                [2541,2583,1592,2674,2323],
                [1212,5214,2325,4235,2519],
            ]
    };
    new Chartist.Bar('#chart-top-products', dataStackedBar, {
            height: "255px",
            stackBars: true,
            axisX: {
                showGrid: false
            },
            axisY: {
                labelInterpolationFnc: function(value) {
                    return (value / 1000) + 'k';
                }
            },
            plugins: [
                Chartist.plugins.tooltip({
                    appendToBody: true
                }),
                Chartist.plugins.legend({
                    legendNames: ['Mobile', 'Laptop', 'Computer']
                })
            ]
    }).on('draw', function(data) {
            if (data.type === 'bar') {
                data.element.attr({
                    style: 'stroke-width: 35px'
                });
            }
    });


    // notification popup
    toastr.options.closeButton = true;
    toastr.options.positionClass = 'toast-bottom-right';
    toastr.options.showDuration = 1000;
    toastr['info']('Hello, welcome to HexaBit, a unique admin Template.');

});

// Customer Overview
function MorrisArea() {
    Morris.Area({
        element: 'Sales_Overview',
        data: [{
                period: '2012',
                SiteA: 0,
                SiteB: 10,

            }, {
                period: '2013',
                SiteA: 106,
                SiteB: 71,

            }, {
                period: '2014',
                SiteA: 68,
                SiteB: 41,

            }, {
                period: '2015',
                SiteA: 89,
                SiteB: 285,

            }, {
                period: '2016',
                SiteA: 185,
                SiteB: 104,

            }, {
                period: '2017',
                SiteA: 146,
                SiteB: 102  ,

            }
        ],
        xkey: 'period',
        ykeys: ['SiteA', 'SiteB'],
        labels: ['Site A', 'Site B'],
        pointSize: 2,
        fillOpacity: 0.9,
        pointStrokeColors: ['#182973', '#29bd73'],
        behaveLikeLine: true,
        gridLineColor: '#e8e8e8',
        lineWidth: 1,
        smooth: true,
        hideHover: 'auto',
        lineColors: ['#182973', '#29bd73'],
        resize: true

    });
}