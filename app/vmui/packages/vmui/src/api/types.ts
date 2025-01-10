export interface MetricBase {
  group: number;
  metric: {
    [key: string]: string;
  };
}

export interface MetricResult extends MetricBase {
  values: [number, string][]
}


export interface InstantMetricResult extends MetricBase {
  value?: [number, string]
  values?: [number, string][]
}

export interface ExportMetricResult extends MetricBase {
  values: number[];
  timestamps: number[];
}

export interface TracingData {
  message: string;
  duration_msec: number;
  children: TracingData[];
}

export interface QueryStats {
  seriesFetched?: string;
  executionTimeMsec?: number;
  resultLength?: number;
  isPartial?: boolean;
}

export interface Logs {
  _msg: string;
  _stream: string;
  _time: string;
  [key: string]: string;
}

export interface LogHits {
  timestamps: string[];
  values: number[];
  total?: number;
  fields: {
    [key: string]: string;
  };
}

export interface ReportMetaData {
  id: number;
  title: string;
  endpoint: string;
  comment: string;
  params: Record<string, string>;
}

export interface LogsFiledValues {
  value: string;
  hits: number;
}
