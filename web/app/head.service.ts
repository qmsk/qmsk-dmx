import { Injectable } from '@angular/core';
import { Http, Headers, RequestOptions } from '@angular/http';

import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/map';

import { Head, Channel } from './head';

@Injectable()
export class HeadService {
  public active: Head = null;

  select(head: Head) {
    console.log("Select head", head);
    this.active = head;
  }
  selected(head: Head): boolean {
    return this.active == head;
  }

  constructor(private http: Http) {

  }

  private post(url, params): Observable<Object> {
    let headers = new Headers({
      'Content-Type': 'application/json',
    });
    let options = new RequestOptions({ headers: headers });

    return this.http.post(url, params, options)
      .map(response => response.json())
    ;
  }

  private decodeHead(headData:Object): Head {
    return Object.assign(new Head, headData, {
      Channels: headData['Channels'].map(channelData => this.decodeChannel(channelData)),
    });
  }
  private decodeChannel(channelData:Object): Channel {
    return Object.assign(new Channel, channelData);
  }

  list(): Observable<Head[]> {
    return this.http.get('/api/heads/')
      .map(response => response.json().map(headData => this.decodeHead(headData)))
    ;
  }

  setHeadChannel(head:Head, channel:Channel, params:Object): Observable<Channel> {
    return this.post(`/api/heads/${head.ID}/channels/${channel.ID}`, params)
      .map(channelParams => Object.assign(channel, channelParams));
  }
}
