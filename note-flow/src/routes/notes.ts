import express from "express";
import {
  createNote,
  getNoteById,
  getNotes,
  renderNoteAsHTML,
} from "../controllers/notesController";

const router = express.Router();

router.post("/", createNote);
router.get("/", getNotes);
router.get("/:id", getNoteById);
router.get("/:id/render", renderNoteAsHTML);

export default router;
