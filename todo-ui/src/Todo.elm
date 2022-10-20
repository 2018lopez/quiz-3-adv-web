-- Filename: src/Todo.elm
module Todo exposing (main)

import Browser
import Html exposing (..)
import Html.Attributes exposing (..)
import Html.Events exposing (onClick, onInput, onSubmit)

-- Model
type alias Model =
    { fieldTitle : String
    ,fieldDetail : String
    , uid : Int
    , todos : List Todo
    }


type alias Todo =
    { id : Int
    , title : String
    , description : String
    , isComplete : Bool
    }
-- Msg
type Msg
    = Add
    | SetFieldTitle String
    | SetFieldDetail String
    | Delete Int
    | CompleteTodo Int Bool

-- Initial Model
initialModel : Model
initialModel =
    { fieldTitle = ""
    , fieldDetail = ""
    , uid = 0
    , todos = []
    }

-- View
view : Model -> Html Msg
view model =
    div [ class "container" ] [ div [ class "todo-container text-left p-24 bg-white shadow-sm rounded flex flex-col mx-auto my-48" ]
        [ header [ ]
            [ h1 [ class "header"] [ text "Todo App" ]
            ]
        , Html.form [ class "formData" ,onSubmit Add ] [
            input
                [ class "todo-input"
                , onInput
                    (\string -> SetFieldTitle string)
                , placeholder "Title"
                , value model.fieldTitle
                ]
                []
            ,textarea
                [ class "todo-textarea"
                ,cols 40, rows 10
                , onInput
                    (\string -> SetFieldDetail string)
                , placeholder "Description"
                , value model.fieldDetail
                ]
                []
            , button [ class "btn", type_ "submit", disabled (model.fieldTitle == "" && model.fieldDetail == "") ] [ text "Add Todo" ]
        ]
        , ul [ class "" ] (List.map viewTask model.todos)
    ]
    
    ]

-- Display task input by user
viewTask : Todo -> Html Msg
viewTask todo =
    li [ class "todo-item-group"  ]
    
        [ 
        div [] [
            div [ class "textTitle"] [
                span [ classList[("completed", todo.isComplete)] ] [  text todo.title ] 
            ]   
            , span [] [ text " : " ] 
            , div [ class "textDescription" ] [
                 span [ classList[("completed", todo.isComplete)] ] [ text todo.description ]
            ]
        ]
        ,div [ class "actionBtn"][
          button
            [ class "btnCompleted", onClick (CompleteTodo todo.id todo.isComplete) ]
            [ text "Completed"]
        , button
            [ class "todo-item-btn", onClick (Delete todo.id)]
            [ text "X" ]
        ]

        ]
-- Update
update : Msg -> Model -> Model
update msg model =
    case msg of
        Add ->
            { model | todos = { id = model.uid, title = model.fieldTitle,  description = model.fieldDetail, isComplete = False } :: model.todos, fieldTitle = "", fieldDetail = "", uid = model.uid + 1 }
        SetFieldTitle str ->
            { model | fieldTitle = str }
        SetFieldDetail str -> 
            { model | fieldDetail = str }
        Delete id ->
            { model | todos = List.filter(\todo -> todo.id /= id) model.todos }
        CompleteTodo id complete ->
            let
                updateTodo todo =
                    if todo.id == id then
                        { todo | isComplete = not complete }
                    else
                        todo
            in
            { model | todos = List.map updateTodo model.todos }


--Main
main : Program () Model Msg
main =
    Browser.sandbox
        { view = view
        , update = update
        , init = initialModel
        }