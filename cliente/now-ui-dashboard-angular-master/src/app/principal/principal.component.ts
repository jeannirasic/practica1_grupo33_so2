import { ServicioService } from '../servicios/servicio.service';
import { ProcesoArbol } from './../interfaces';
import {MatPaginator, MatSort, MatTableDataSource} from '@angular/material';
import { Resumen, Procesos } from '../interfaces';
import * as go from 'gojs';
import { Component, OnInit, Output, Input, EventEmitter, ChangeDetectorRef, ViewChild } from '@angular/core';
import { Subscription, Observable, timer } from 'rxjs';
import * as moment from 'moment';

let $ = go.GraphObject.make;

@Component({
  selector: 'app-principal',
  templateUrl: './principal.component.html',
  styleUrls: ['./principal.component.scss']
})
export class PrincipalComponent implements OnInit {

  private subscription: Subscription;
  @Output() TimerExpired: EventEmitter<any> = new EventEmitter<any>();
  @Input() SearchDate: moment.Moment = moment();
  @Input() ElapsTime = 5;
  searchEndDate: moment.Moment;
  remainingTime: number;
  minutes: number;
  seconds: number;
  everySecond: Observable<number> = timer(0, 1000);

  public myDiagram: go.Diagram = null;

  resumen: Resumen = {total: 0, ejecucion: 0, suspendidos: 0, detenidos: 0, zombie: 0};
  listaProcesos: Procesos[] = new Array();
  nodeDataArray: ProcesoArbol[] = new Array();
  displayedColumns: string[] = ['pid', 'nombre', 'usuario', 'estado', 'ram', 'accion'];
  dataSource: MatTableDataSource<Procesos>;
  @ViewChild(MatPaginator) paginator: MatPaginator;
  @ViewChild(MatSort) sort: MatSort;

  constructor(private servicio: ServicioService, private ref: ChangeDetectorRef) {
    this.actualizarDatos();
    this.searchEndDate = this.SearchDate.add(this.ElapsTime, 'seconds');
  }

  ngOnInit() {
    this.subscription = this.everySecond.subscribe((seconds) => {
      const currentTime: moment.Moment = moment();
      this.remainingTime = this.searchEndDate.diff(currentTime);
      this.remainingTime = this.remainingTime / 1000;
      if (this.remainingTime <= 0) {
        this.SearchDate = moment();
        this.searchEndDate = this.SearchDate.add(this.ElapsTime, 'seconds');
        this.TimerExpired.emit();
        console.log('Se acabo');
        this.actualizarDatos();
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

  actualizarDatos() {
    this.servicio.informacionPrincipal().subscribe(data => {
      for (let i = 0; i < data.length; i++) {
        if (data[i].estado === 'T') {
          data[i].booleano = false;
        } else {
          data[i].booleano = true;
        }
      }
      this.listaProcesos = data;
      this.dataSource = new MatTableDataSource(this.listaProcesos);
      this.dataSource.paginator = this.paginator;
      this.dataSource.sort = this.sort;

      // Creo el arbol
      this.nodeDataArray = new Array();
      for (let i = 0; i < data.length; i++) {
        const nuevo: ProcesoArbol = {
          key: data[i].pid,
          name: 'PID: ' + data[i].pid + ', \n Nombre: ' + data[i].nombre,
          parent: data[i].ppid
        };
        this.nodeDataArray.push(nuevo);
      }
      this.myDiagram.model = new go.TreeModel(this.nodeDataArray);

      // Reviso lo de las estadisticas generales
      let contadorEjecucion = 0, contadorSuspendidos = 0, contadorDetenidos = 0, contadorZombie = 0;
      for (let i = 0; i < data.length; i++) {
        if (data[i].estado === 'R') {
          contadorEjecucion++;
        } else if (data[i].estado === 'S') {
          contadorSuspendidos++;
        } else if (data[i].estado === 'T') {
          contadorDetenidos++;
        } else if (data[i].estado === 'Z') {
          contadorZombie++;
        }
      }

      const resumenTemp: Resumen = {
        total: data.length,
        ejecucion: contadorEjecucion,
        suspendidos: contadorSuspendidos,
        detenidos: contadorDetenidos,
        zombie: contadorZombie
      };
      this.resumen = resumenTemp;
    }, error => {
      alert('Ha ocurrido un error al obtener la lista de procesos');
    });

  }

  terminar(e) {
    this.servicio.matarProceso(e.nombre).subscribe(data => {
      alert('El proceso se ha eliminado existosamente');
    }, error => {
      alert('Hubo un error al eliminar el proceso');
    });
  }

  // tslint:disable-next-line:use-life-cycle-interface
  ngAfterViewInit() {
    this.myDiagram =
        $(go.Diagram, 'myDiagramDiv',
          {
            'toolManager.hoverDelay': 100,
            allowCopy: false,
            layout:
              $(go.TreeLayout,
                { angle: 90, nodeSpacing: 10, layerSpacing: 40, layerStyle: go.TreeLayout.LayerUniform })
          });

      this.myDiagram.add(
        $(go.Part, 'Table',
          { position: new go.Point(300, 10), selectable: false }
        ));


      this.myDiagram.nodeTemplate =
        $(go.Node, 'Auto',
          { deletable: false },
          new go.Binding('text', 'name'),
          $(go.Shape, 'Rectangle',
            {
              fill: '#E16D6D',
              stroke: null, strokeWidth: 0,
              stretch: go.GraphObject.Fill,
              alignment: go.Spot.Center
            }),
          $(go.TextBlock,
            {
              font: '700 12px Droid Serif, sans-serif',
              textAlign: 'center',
              margin: 10, maxSize: new go.Size(80, NaN)
            },
            new go.Binding('text', 'name'))
        );
      this.myDiagram.linkTemplate =
        $(go.Link,
          { routing: go.Link.Orthogonal, corner: 5, selectable: false },
          $(go.Shape, { strokeWidth: 3, stroke: '#424242' }));


      this.myDiagram.model = new go.TreeModel(this.nodeDataArray);
  }

  applyFilter(filterValue: string) {
    this.dataSource.filter = filterValue.trim().toLowerCase();

    if (this.dataSource.paginator) {
      this.dataSource.paginator.firstPage();
    }
  }
}
