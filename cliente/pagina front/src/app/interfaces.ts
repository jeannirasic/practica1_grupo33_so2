export interface Resumen {
    ejecucion: number;
    suspendidos: number;
    detenidos: number;
    zombie: number;
    total: number;
}

export interface Procesos {
    PID: number;
    Nombre: string;
    Usuario: string;
    Estado: string;
    PorcentajeRam: number;
    PPID: number;
    booleano?: boolean;
}

export interface ProcesoArbol {
    key: number;
    name: string;
    parent?: number;
}

export interface Ram {
    Memoria_Total: number;
    Memoria_en_uso: number;
    Porcentaje_en_uso: number;
}

