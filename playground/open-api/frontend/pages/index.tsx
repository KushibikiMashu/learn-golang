import aspida from "@aspida/fetch"
import api from '../api/$api'
import {Item} from "../api/@types";
import useSWR, {mutate} from "swr";
import {Methods} from "../api";

const url = 'http://localhost:3002/'

type Todo = Required<Item>

const client = api(aspida(fetch, {baseURL: url}))

type ApiQuery = Methods['get']['query']
type ApiReqBody = Methods['post']['reqBody']

const fetchAllTodo = async(query: ApiQuery)=> {
  return await client.$get({query})
}

const useFetchAllTodos = () => {
  const {data = [], error} = useSWR<Todo[], Error>('/', () => fetchAllTodo({since: 0, limit: 20}))
  const todos = data.length === 0 ? data : data.sort((prev, next) => prev.id > next.id ? 1: -1)

  return {todos, error}
}

const postTodo = async (body: ApiReqBody) => {
  return await client.$post({body})
}

const addTodo = async (todos: Todo[]) => {
  const description = Math.random().toString(36).slice(6)
  const newTodo = {description, completed: false}
  const newId = todos.length === 0 ? 1: todos[todos.length - 1].id + 1

  await mutate('/', [...todos, {...newTodo, id: newId}])
  await postTodo(newTodo)
  await mutate('/')
}

const putTodoComplete = async (todos: Todo[], newTodo: Todo) => {
  const id = newTodo.id

  await mutate('/', todos.map(todo => todo.id === id ? newTodo : todo))
  await client._id(id).put({body: {description: newTodo.description, completed: newTodo.completed}})
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
  const handleChange = async (todo: Todo) => await putTodoComplete(todos, todo)

  if (!todos) {
    return <div>loading</div>
  }

  if (error) {
    return <div>error: {String(error)}</div>
  }

  const AddButton = () => <button style={{marginLeft: 40}} onClick={handleAddClick}>Add Todo</button>

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
        <li style={{display: 'flex', justifyContent: 'space-between', width: 200, margin: '4px 0'}} key={i}>
          <span style={{width: 140}}>{todo.id}: {todo.description}</span>
          <input type="checkbox" checked={!!todo.completed} onChange={() => handleChange({...todo, completed: !todo.completed})} />
          <button style={{marginLeft: 8}} onClick={() => handleDeleteClick(todo.id)}>??</button>
        </li>
      ))}
    </ul>
      <AddButton />
    </>
  )
}
