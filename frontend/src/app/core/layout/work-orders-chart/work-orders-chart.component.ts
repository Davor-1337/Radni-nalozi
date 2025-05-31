import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { Chart } from 'chart.js/auto';
import { WorkOrderService } from '../../services/workOrder.service';

interface WorkOrderStat {
  MonthNumber: number;
  Count: number;
}

@Component({
  selector: 'app-work-orders-chart',
  templateUrl: './work-orders-chart.component.html',
  styleUrls: ['./work-orders-chart.component.scss'],
})
export class WorkOrdersChartComponent implements OnInit {
  @ViewChild('workOrdersChart', { static: true }) chartRef!: ElementRef;
  private chart!: Chart;
  isLoading = true;

  constructor(private workOrderService: WorkOrderService) {}

  ngOnInit(): void {
    this.loadChartData();
  }

  loadChartData(): void {
    this.workOrderService.getWorkOrderStats().subscribe({
      next: (stats) => {
        this.createChart(this.prepareChartData(stats));
        this.isLoading = false;
        console.log('API Response:', stats);
      },
      error: (err) => {
        console.error('Error fetching work order stats:', err);
        this.isLoading = false;
        this.createChart(this.prepareChartData([]));
      },
    });
  }

  prepareChartData(stats: WorkOrderStat[]): { data: number[] } {
    const monthlyData = new Array(12).fill(0);

    stats.forEach((stat) => {
      monthlyData[stat.MonthNumber - 1] = stat.Count;
    });

    return { data: monthlyData };
  }

  createChart(chartData: { data: number[] }): void {
    const canvas = this.chartRef.nativeElement as HTMLCanvasElement;

    if (this.chart) {
      this.chart.destroy();
    }

    this.chart = new Chart(canvas, {
      type: 'line',
      data: {
        labels: [
          'Jan',
          'Feb',
          'Mar',
          'Apr',
          'May',
          'Jun',
          'Jul',
          'Aug',
          'Sep',
          'Oct',
          'Nov',
          'Dec',
        ],
        datasets: [
          {
            label: 'Work orders by months',
            data: chartData.data,
            borderColor: 'rgba(54, 162, 235, 1)',
            backgroundColor: 'rgba(54, 162, 235, 0.2)',
            tension: 0.4,
            fill: true,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,
        plugins: {
          legend: {
            display: false,
            position: 'top',
          },
        },
        scales: {
          x: {
            ticks: {
              color: '#515961',
              font: {
                size: 12,
                family: 'Inter',
              },
            },
            title: { display: false, text: 'Months' },
          },
          y: {
            beginAtZero: true,
            ticks: {
              stepSize: 2,
              color: '#515961',
              font: {
                size: 12,
                family: 'Inter',
              },
            },
            title: {
              display: true,
              text: 'Completed work orders',
              color: '#041226',
              font: {
                size: 12,
                family: 'Inter',
                weight: 'normal',
              },
            },
          },
        },
      },
    });
  }
}
