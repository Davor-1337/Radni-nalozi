import { Component, ElementRef, OnInit, ViewChild } from '@angular/core';
import { Chart } from 'chart.js/auto';
import { MaterialServiceService } from '../../../core/services/material-service.service';

@Component({
  selector: 'app-chart',
  templateUrl: './chart.component.html',
  styleUrls: ['./chart.component.scss'],
})
export class ChartComponent implements OnInit {
  @ViewChild('materialsChart', { static: true }) chartRef!: ElementRef;
  private chart!: Chart;

  constructor(private materialsService: MaterialServiceService) {}

  ngOnInit(): void {
    this.loadChartData();
  }

  loadChartData(): void {
    this.materialsService.getMaterialsStats().subscribe({
      next: (stats) => {
        this.createChart(stats);
      },
      error: (err) => {
        console.error('Error loading material stats:', err);
      },
    });
  }

  createChart(stats: any[]): void {
    const canvas = this.chartRef.nativeElement as HTMLCanvasElement;

    if (this.chart) {
      this.chart.destroy();
    }

    const labels = stats.map((item) => item.kategorija);
    const data = stats.map((item) => item.ukupno_utroseno);

    const backgroundColors = [
      'rgba(16, 52, 94, 0.8)',
      'rgba(88, 102, 110, 0.8)',
      'rgba(134, 38, 51, 0.8)',
      'rgba(65, 45, 84, 0.8)',
      'rgba(128, 90, 50, 0.8)',
    ];

    this.chart = new Chart(canvas, {
      type: 'pie',
      data: {
        labels: labels,
        datasets: [
          {
            label: 'Materijali po kategorijama',
            data: data,
            backgroundColor: backgroundColors,
            borderWidth: 1,
          },
        ],
      },
      options: {
        responsive: true,
        maintainAspectRatio: false,

        plugins: {
          legend: {
            position: 'bottom',

            labels: {
              color: '#515961',
              padding: 40,
              font: {
                size: 12,
                family: 'Inter',
              },
            },
          },
          tooltip: {
            callbacks: {
              label: (context) => {
                const label = context.label || '';
                const value = context.raw as number;
                const total = context.dataset.data.reduce(
                  (a: number, b: number) => a + b,
                  0
                );
                const percentage = Math.round((value / total) * 100);
                return `${label}: ${value} kom (${percentage}%)`;
              },
            },
          },
        },
      },
    });
  }
}
