import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Resumen, Procesos, Ram, Cpu } from '../interfaces';

@Injectable({
  providedIn: 'root'
})
export class ServicioService {

  url = 'http://localhost:3000/';

  constructor(private http: HttpClient) { }

  // La informacion de los procesos
  informacionPrincipal() {
    return this.http.get<Procesos[]>(this.url + 'principal');
  }

  // Devuelve los datos de la RAM
  informacionRam() {
    return this.http.get<Ram>(this.url + 'ram');
  }

  // Devuelve la informacion del CPU
  informacionCpu() {
    return this.http.get<Cpu>(this.url + 'cpu');
  }

  // Envia el pid del proceso para matarlo
  matarProceso(pid: string) {
    return this.http.post<any>(this.url + 'kill/' + pid, '');
  }
}
