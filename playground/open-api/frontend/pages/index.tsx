import aspida from "@aspida/fetch"
import api from '../api/$api'
import {Item} from "../api/@types";
import useSWR, {mutate} from "swr";
import {Methods} from "../api";

const url = 'http://127.0.0.1:3002/'

type Todo = Required<Item>

const client = api(aspida(fetch, {baseURL: url}))

type ApiQuery = Methods['get']['query']
type ApiReqBody = Methods['post']['reqBody']

const fetchAllTodo = async(query: ApiQuery)=> {
  return await client.$get({query})
}

const useFetchAllTodos = () => {
  const {data = [], error} = useSWR<Todo[], Error>('/', () => fetchAllTodo({since: 0, limit: 20}))

  return {todos: data, error}
}

const postTodo = async (body: ApiReqBody) => {
  return await client.$post({body})
}

const addTodo = async (todos: Todo[]) => {
  const description = Math.random().toString(36).slice(6)
  const newTodo = {description, completed: false}

  await mutate('/', [...todos, {...newTodo, id: todos.length + 1}])
  await postTodo(newTodo)
  await mutate('/')
}

const deleteTodo = async (id: number, todos: Item[]) => {
  await client._id(id).delete()

  const newTodos = todos.filter(todo => todo.id !== id)
  await mutate('/', newTodos)
  await mutate('/')
}

export default function Index() {
  const {todos, error} = useFetchAllTodos()
  const handleAddClick = async () => await addTodo(todos)
  const handleDeleteClick = async (id: number) => await deleteTodo(id, todos)

  if (!todos) {
    return <div>loading</div>
  }

  if (error) {
    return <div>error: {String(error)}</div>
  }

  const AddButton = () => <button onClick={handleAddClick}>Add Todo</button>

  if (todos.length === 0) {
    return (
      <>
        <div>empty</div>
        <AddButton />
      </>
    )
  }

  return (
    <>
    <ul>
      {todos.map((todo: Todo, i) => (
        <li style={{display: 'flex', justifyContent: 'space-between', width: 180, margin: '4px 0'}} key={i}>
          {todo.id}: {todo.description}
          <button style={{marginLeft: 8}} onClick={() => handleDeleteClick(todo.id)}>×</button>
        </li>
      ))}
    </ul>
      <AddButton />
    </>
  )
}
