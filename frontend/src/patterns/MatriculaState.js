class EstadoMatricula {
  constructor(nombre, siguientesEstados) {
    this.nombre = nombre;
    this.siguientesEstados = siguientesEstados;
  }

  acciones() {
    return this.siguientesEstados;
  }

  permite(destino) {
    return this.siguientesEstados.includes(destino);
  }
}

export class EstadoActiva extends EstadoMatricula {
  constructor() {
    super('Activa', ['Finalizada', 'Anulada']);
  }
}

export class EstadoFinalizada extends EstadoMatricula {
  constructor() {
    super('Finalizada', []);
  }
}

export class EstadoAnulada extends EstadoMatricula {
  constructor() {
    super('Anulada', []);
  }
}

export class EstadoPendiente extends EstadoMatricula {
  constructor() {
    super('Pendiente', ['Activa', 'Anulada']);
  }
}

const estados = {
  Activa: new EstadoActiva(),
  Finalizada: new EstadoFinalizada(),
  Anulada: new EstadoAnulada(),
  Pendiente: new EstadoPendiente(),
};

// El objeto de estado define las acciones válidas; la vista no duplica las reglas.
export function obtenerEstadoMatricula(nombre) {
  return estados[nombre] ?? estados.Activa;
}
