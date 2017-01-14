import { DMX, Value, Color, ChannelType, HeadType, HeadConfig } from './types';

// GET
export interface APIChannel {
  ID:       string;
  Type:     ChannelType;
  Address:  number;
  DMX:      DMX;
  Value:    Value;
}

export interface APIIntensity {
  Intensity: Value;
}

export interface APIColor extends Color {

}

export interface APIHead {
  ID:       string;
  Type:     HeadType;
  Config:   HeadConfig;

  Channels:   {[id: string]: APIChannel};
  Intensity?: APIIntensity;
  Color?:     APIColor;
}
export type APIHeads = {[id: string]: APIHead};

export interface APIGroup {
  ID:     string;
  Heads:  string[];

  Intensity?: APIIntensity;
  Color?:     APIColor;
}
export type APIGroups = {[id: string]: APIGroup};


export interface API {
  Heads:  APIHeads;
  Groups: APIGroups;
}

// POST
export interface APIChannelParameters {
  DMX?: DMX;
  Value?: Value;
}

export interface APIParameters {
  Intensity?: APIIntensity;
  Color?:     APIColor;
}

export interface APIHeadParameters extends APIParameters {
  Channels?:  {[ID: string]: APIChannelParameters};
};

// WebSocket
export interface APIEvents extends API {
  
}