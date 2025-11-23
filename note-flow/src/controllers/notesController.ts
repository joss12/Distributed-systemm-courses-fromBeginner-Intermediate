import { Request, Response } from "express";
import { Note } from "../models/Note";
import { marked } from "marked";

let notes: Note[] = [];

export const createNote = (req: Request, res: Response) => {
  const { title, content } = req.body;
  if (!title || !content)
    return res.status(400).json({ message: "Missing title or content" });

  const newNote: Note = {
    id: Date.now().toString(),
    title,
    content,
    createdAt: new Date(),
    updatedAt: new Date(),
  };

  notes.push(newNote);
  res.status(201).json(newNote);
};

export const getNotes = (req: Request, res: Response) => {
  res.json(notes);
};

export const getNoteById = (req: Request, res: Response) => {
  const note = notes.find((n) => n.id === req.params.id);
  if (!note) return res.status(404).json({ message: "Note not found" });
  res.json(note);
};

export const renderNoteAsHTML = (req: Request, res: Response) => {
  const note = notes.find((n) => n.id === req.params.id);
  if (!note) return res.status(404).json({ message: "Note not found" });
  const html = marked(note.content);
  res.send(`<h1>${note.title}</h1><article>${html}</article>`);
};
