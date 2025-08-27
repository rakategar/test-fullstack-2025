// Raka Tegar Wicaksono
// tegarrakaw@gmail.com
// Fullstack

// 1
function f(n: number): number {
  const factorial = (x: number): number => (x <= 1 ? 1 : x * factorial(x - 1));

  const result = factorial(n) / Math.pow(2, n);
  return Math.ceil(result);
}
