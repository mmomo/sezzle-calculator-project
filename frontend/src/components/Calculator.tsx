import { useState, type FormEvent } from 'react'
import { calculate, CalculatorApiError, type Operation } from '../api/calculator'
import './Calculator.css'

const OPERATIONS: { value: Operation; label: string; unary?: boolean }[] = [
  { value: 'add', label: 'Add (a + b)' },
  { value: 'subtract', label: 'Subtract (a − b)' },
  { value: 'multiply', label: 'Multiply (a × b)' },
  { value: 'divide', label: 'Divide (a ÷ b)' },
  { value: 'exponentiate', label: 'Exponentiate (a ^ b)' },
  { value: 'sqrt', label: 'Square root (√a)', unary: true },
  { value: 'percentage', label: 'Percentage (a% of b)' },
]

function isUnary(operation: Operation): boolean {
  return OPERATIONS.find((candidate) => candidate.value === operation)?.unary ?? false
}

export function Calculator() {
  const [operation, setOperation] = useState<Operation>('add')
  const [a, setA] = useState('')
  const [b, setB] = useState('')
  const [result, setResult] = useState<number | null>(null)
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(false)

  const unary = isUnary(operation)
  const aValid = a.trim() !== '' && !Number.isNaN(Number(a))
  const bValid = unary || (b.trim() !== '' && !Number.isNaN(Number(b)))
  const canSubmit = aValid && bValid && !loading

  async function handleSubmit(event: FormEvent) {
    event.preventDefault()
    if (!canSubmit) return

    setLoading(true)
    setError(null)
    setResult(null)

    try {
      const value = await calculate({
        operation,
        a: Number(a),
        b: unary ? undefined : Number(b),
      })
      setResult(value)
    } catch (err) {
      setError(
        err instanceof CalculatorApiError
          ? err.message
          : 'Something went wrong. Please try again.',
      )
    } finally {
      setLoading(false)
    }
  }

  return (
    <form className="calculator" onSubmit={handleSubmit}>
      <h1>Calculator</h1>

      <label className="field">
        <span>Operation</span>
        <select
          value={operation}
          onChange={(event) => setOperation(event.target.value as Operation)}
        >
          {OPERATIONS.map((op) => (
            <option key={op.value} value={op.value}>
              {op.label}
            </option>
          ))}
        </select>
      </label>

      <label className="field">
        <span>a</span>
        <input
          type="text"
          inputMode="decimal"
          value={a}
          onChange={(event) => setA(event.target.value)}
          placeholder="e.g. 4"
        />
      </label>

      {!unary && (
        <label className="field">
          <span>b</span>
          <input
            type="text"
            inputMode="decimal"
            value={b}
            onChange={(event) => setB(event.target.value)}
            placeholder="e.g. 2"
          />
        </label>
      )}

      <button type="submit" disabled={!canSubmit}>
        {loading ? 'Calculating…' : '='}
      </button>

      {result !== null && (
        <p className="result" role="status">
          Result: {result}
        </p>
      )}

      {error !== null && (
        <p className="error" role="alert">
          {error}
        </p>
      )}
    </form>
  )
}
