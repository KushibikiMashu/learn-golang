import aspida from "@aspida/fetch"
import api from '../api/$api'
import {useEffect, useState} from "react";
import {Item} from "../api/@types";
import useSWR from "swr";

const url = 'http://127.0.0.1:3002/'

const client = api(aspida(fetch, {baseURL: url}))

const fetchAllTodo = async(
  {since = 0, limit = 20}:
    {since: number, limit: number}
): Promise<Item[]> => {
  return await client.$get({query: {since, limit}})
}

const useFetchTodos = () => {
  const {data, error} = useSWR<Item[], Error>('/', () => fetchAllTodo({since: 0, limit: 20}))

  return {todos: data, error}
}

export default function Index() {
  const {todos,  error } = useFetchTodos()

  const handleClick = async () => {}

  if (!todos) {
    return <div>loading</div>
  }

  if (error) {
    return <div>error: {String(error)}</div>
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
