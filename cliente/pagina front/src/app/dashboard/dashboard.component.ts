import { Component, OnInit } from '@angular/core';
import * as Chartist from 'chartist';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-dashboard',
  templateUrl: './dashboard.component.html',
  styleUrls: ['./dashboard.component.css']
})
export class DashboardComponent implements OnInit {
  constructor(private route: ActivatedRoute, private router: Router  ) { }

  ngOnInit() {
  }

  procesos() {
    this.router.navigate(['/procesos']);
  }

  memoria() {
    this.router.navigate(['/ram']);
  }
}
