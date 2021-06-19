import aspida from "@aspida/fetch"
import api from '../api/$api'
import {useEffect, useState} from "react";
import {Item} from "../api/@types";

const url = 'http://127.0.0.1:3002/'

const client = api(aspida(fetch, {baseURL: url}))

const fetchAllTodo = async(
  {since = 0, limit = 20}:
    {since: number, limit: number}
) => {
  return await client.$get({query: {since, limit}})
}

const useFetchTodos = () => {
  const [todos, setTodos] = useState<Item[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<Error | null>(null)

  useEffect(() => {
    (async() => {
      setLoading(true)
      try {
        const todos = await fetchAllTodo({since: 0, limit: 20})
        setTodos(todos)
      } catch (e) {
        setError(e)
      } finally {
        setLoading(false)
      }
    })()
  }, [])

  return {todos, loading, error}
}

export default function Index() {
  const {todos, loading, error } = useFetchTodos()
  const handleClick = async () => {}

  if (loading) {
    return <div>loading</div>
  }

  if (error !== null) {
    return <div>error: {error.message}</div>
  }

  if (todos.length === 0) {
    return (
      <>
        <div>empty</div>
        <button onClick={handleClick}>Add Todo</button>
      </>
    )
  }

  return (
    <>
    <ul>
      {todos.map((todo: Item, i) => (<li key={i}>{todo.id}: {todo.description}</li>))}
    </ul>
    </>
  )
}
