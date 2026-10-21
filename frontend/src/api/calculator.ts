export type Operation =
  | 'add'
  | 'subtract'
  | 'multiply'
  | 'divide'
  | 'exponentiate'
  | 'sqrt'
  | 'percentage'

export interface CalculateRequest {
  operation: Operation
  a: number
  b?: number
}

interface CalculateResponse {
  result: number
}

interface ErrorResponse {
  error: string
}

export class CalculatorApiError extends Error {
  status: number

  constructor(message: string, status: number) {
    super(message)
    this.name = 'CalculatorApiError'
    this.status = status
  }
}

export async function calculate(request: CalculateRequest): Promise<number> {
  const response = await fetch('/api/calculate', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify(request),
  })

  if (!response.ok) {
    const body = (await response.json()) as ErrorResponse
    throw new CalculatorApiError(body.error, response.status)
  }

  const body = (await response.json()) as CalculateResponse
  return body.result
}
