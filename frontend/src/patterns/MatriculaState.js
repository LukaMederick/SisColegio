// Patrón State: Define los estados válidos y transiciones permitidas
export const Estados = {
  Pendiente: ['Activa', 'Anulada'],
  Activa: ['Finalizada', 'Anulada'],
  Finalizada: [],
  Anulada: [],
};

export const puedeCambiar = (actual, nuevo) =>
  Estados[actual]?.includes(nuevo) ?? false;
