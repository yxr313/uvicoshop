export type ElementType = "image" | "text" | "shape";

export interface CanvasElement {
    id: string;
    type: ElementType;
    x: number;
    y: number;
    width?: number;
    height?: number;
    rotation: number;
    visible: boolean;
    locked: boolean;
    // Props de Imagen
    src?: string;
    // Props de Texto
    text?: string;
    fontSize?: number;
    fontFamily?: string;
    fill?: string;
}

export interface Design {
    id: string;
    name: string;
    mockupUrl?: string;
    elements: CanvasElement[];
    createdAt: string;
    updatedAt: string;
}