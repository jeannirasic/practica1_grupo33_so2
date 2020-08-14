import { ServicioService } from './../servicios/servicio.service';
import { Ram } from '../interfaces';
import { Component, OnInit, Output, Input, EventEmitter, ChangeDetectorRef } from '@angular/core';
import { Subscription, Observable, timer } from 'rxjs';
import * as moment from 'moment';

@Component({
  selector: 'app-ram',
  templateUrl: './ram.component.html',
  styleUrls: ['./ram.component.scss']
})
export class RamComponent implements OnInit {

  private subscription: Subscription;
  @Output() TimerExpired: EventEmitter<any> = new EventEmitter<any>();
  @Input() SearchDate: moment.Moment = moment();
  @Input() ElapsTime = 3;
  searchEndDate: moment.Moment;
  remainingTime: number;
  minutes: number;
  seconds: number;
  everySecond: Observable<number> = timer(0, 1000);

  ram: Ram;
  public lineBigDashboardChartType;
  public gradientStroke;
  public chartColor;
  public canvas: any;
  public ctx;
  public gradientFill;
  public lineBigDashboardChartData: Array<any>;
  public lineBigDashboardChartOptions: any;
  public lineBigDashboardChartLabels: Array<any>;
  public lineBigDashboardChartColors: Array<any>;


  datosGrafica: number[] = new Array();
  labelsGrafica: string[] = new Array();

  public chartClicked(e: any): void {
    console.log(e);
  }

  public chartHovered(e: any): void {
    console.log(e);
  }
  public hexToRGB(hex, alpha) {
    const r = parseInt(hex.slice(1, 3), 16),
      g = parseInt(hex.slice(3, 5), 16),
      b = parseInt(hex.slice(5, 7), 16);

    if (alpha) {
      return 'rgba(' + r + ', ' + g + ', ' + b + ', ' + alpha + ')';
    } else {
      return 'rgb(' + r + ', ' + g + ', ' + b + ')';
    }
  }
  constructor(private servicio: ServicioService, private ref: ChangeDetectorRef) {
    this.actualizar();
    this.searchEndDate = this.SearchDate.add(this.ElapsTime, 'seconds');
  }

  ngOnInit() {
    this.grafica();
    this.subscription = this.everySecond.subscribe((seconds) => {
      const currentTime: moment.Moment = moment();
      this.remainingTime = this.searchEndDate.diff(currentTime);
      this.remainingTime = this.remainingTime / 1000;
      if (this.remainingTime <= 0) {
        this.SearchDate = moment();
        this.searchEndDate = this.SearchDate.add(this.ElapsTime, 'seconds');
        this.TimerExpired.emit();
        console.log('Se acabo');
        this.actualizar();
      } else {
        this.minutes = Math.floor(this.remainingTime / 60);
        this.seconds = Math.floor(this.remainingTime - this.minutes * 60);
      }
      this.ref.markForCheck();
      });
  }

  // tslint:disable-next-line:use-life-cycle-interface
  ngOnDestroy(): void {
  this.subscription.unsubscribe();
  }

  actualizar() {
    this.servicio.informacionRam().subscribe(data => {
      this.ram = data;
      this.datosGrafica.push(data.Porcentaje_en_uso);
      this.labelsGrafica.push(this.timeGenerate());
      this.grafica();
    });
  }


  timeGenerate(): string {
    const time = new Date();
    const hours = time.getHours();
    const minutes = time.getMinutes();
    const seconds = time.getSeconds();

    let cadenaHoras = '', cadenaMinutos = '', cadenaSegundos = '';
    if (hours < 10) {
      cadenaHoras = '0' + hours.toString();
    } else {
      cadenaHoras = hours.toString();
    }

    if (minutes < 10) {
      cadenaMinutos = '0' + minutes.toString();
    } else {
      cadenaMinutos = minutes.toString();
    }

    if (seconds < 10) {
      cadenaSegundos = '0' + seconds.toString();
    } else {
      cadenaSegundos = seconds.toString();
    }

    const cadenaFinal = cadenaHoras + ':' + cadenaMinutos + ':' + cadenaSegundos;
    return cadenaFinal;
}


  grafica() {
    this.chartColor = '#FFFFFF';
    this.canvas = document.getElementById('bigDashboardChart');
    this.ctx = this.canvas.getContext('2d');

    this.gradientStroke = this.ctx.createLinearGradient(500, 0, 100, 0);
    this.gradientStroke.addColorStop(0, '#80b6f4');
    this.gradientStroke.addColorStop(1, this.chartColor);

    this.gradientFill = this.ctx.createLinearGradient(0, 200, 0, 50);
    this.gradientFill.addColorStop(0, 'rgba(128, 182, 244, 0)');
    this.gradientFill.addColorStop(1, 'rgba(255, 255, 255, 0.24)');

    this.lineBigDashboardChartData = [
        {
          label: 'Data',

          pointBorderWidth: 1,
          pointHoverRadius: 7,
          pointHoverBorderWidth: 2,
          pointRadius: 5,
          fill: true,
          lineTension: 0,
          borderWidth: 2,
          data: this.datosGrafica
        }
      ];
      this.lineBigDashboardChartColors = [
       {
         backgroundColor: this.gradientFill,
         borderColor: this.chartColor,
         pointBorderColor: this.chartColor,
         pointBackgroundColor: '#2c2c2c',
         pointHoverBackgroundColor: '#2c2c2c',
         pointHoverBorderColor: this.chartColor,
       }
     ];
    this.lineBigDashboardChartLabels = this.labelsGrafica;
    this.lineBigDashboardChartOptions = {

          layout: {
              padding: {
                  left: 20,
                  right: 20,
                  top: 0,
                  bottom: 0
              }
          },
          maintainAspectRatio: false,
          tooltips: {
            backgroundColor: '#fff',
            titleFontColor: '#333',
            bodyFontColor: '#666',
            bodySpacing: 4,
            xPadding: 12,
            mode: 'nearest',
            intersect: 0,
            position: 'nearest'
          },
          legend: {
              position: 'bottom',
              fillStyle: '#FFF',
              display: false
          },
          scales: {
              yAxes: [{
                  ticks: {
                      fontColor: 'rgba(255,255,255,0.4)',
                      fontStyle: 'bold',
                      beginAtZero: true,
                      maxTicksLimit: 5,
                      padding: 10
                  },
                  gridLines: {
                      drawTicks: true,
                      drawBorder: false,
                      display: true,
                      color: 'rgba(255,255,255,0.1)',
                      zeroLineColor: 'transparent'
                  }

              }],
              xAxes: [{
                  gridLines: {
                      zeroLineColor: 'transparent',
                      display: false,

                  },
                  ticks: {
                      padding: 10,
                      fontColor: 'rgba(255,255,255,0.4)',
                      fontStyle: 'bold'
                  }
              }]
          }
    };

    this.lineBigDashboardChartType = 'line';

  }

}
