export type ViewMode = 'table' | 'grid' | 'card' | 'split' | 'miller';

export interface ViewModeOption {
  id: ViewMode;
  label: string;
  icon: string;
}

export type StatusType = 'success' | 'error' | 'info' | 'warning';
