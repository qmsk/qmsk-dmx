import { Component, Input, HostBinding } from '@angular/core';

import { Value, Color, Colors } from './types';
import { Head, Group, ColorParameter } from './head';
import { APIService } from './api.service';

@Component({
  moduleId: module.id,
  selector: 'dmx-color',
  host: { class: 'view split' },

  templateUrl: 'color.component.html',
  styleUrls: [ 'color.component.css' ],
})
export class ColorComponent {
  colors: Colors;
  color: Color;
  heads: Set<Head>;
  groups: Set<Group>;

  constructor (private api: APIService) {
    this.heads = new Set<Head>();
    this.groups = new Set<Group>();
  }

  listHeads(): Head[] {
    return this.api.listHeads(head => head.ID, head => !!head.Color);
  }
  listGroups(): Group[] {
    return this.api.listGroups(group => group.ID, group => !!group.Color);
  }

  headActive(head: Head): boolean {
    return this.heads.has(head);
  }
  groupActive(group: Group): boolean {
    return this.groups.has(group);
  }

  /* Build new colors map from active heads */
  loadColors(): Colors {
    // XXX: just return from first selected group or head...
    // TODO: merge color maps from multiple heads?
    for (let group of Array.from(this.groups)) {
      return group.Colors;
    }

    for (let head of Array.from(this.heads)) {
      return head.Colors;
    }

    return null;
  }

  selectHead(head: Head) {
    this.heads.add(head);

    this.colors = this.loadColors();
    this.color = head.Color;
  }
  selectGroup(group: Group) {
    this.groups.add(group);

    this.colors = this.loadColors();
    this.color = group.Color;
  }

  unselect() {
    if (this.colors = this.loadColors()) {
      // keep selected color
    } else {
      this.color = null;
    }
  }

  unselectHead(head: Head) {
    this.heads.delete(head);
    this.unselect();
  }
  unselectGroup(group: Group) {
    this.groups.delete(group);

    this.unselect();
  }

  apply(color: Color) {
    this.heads.forEach((head) => {
      head.Color.apply(color);
    });
    this.groups.forEach((group) => {
      group.Color.apply(color);
    });
  }
}
