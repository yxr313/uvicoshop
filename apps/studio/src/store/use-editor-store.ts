import { create } from "zustand";
import { CanvasElement } from "@/types";

interface EditorState {
    elements: CanvasElement[];
    selectedId: string | null;
    activeTool: "text" | "image" | "shapes" | "ai" | null;

    // Acciones
    setElements: (elements: CanvasElement[]) => void;
    addElement: (element: CanvasElement) => void;
    updateElement: (id: string, newProps: Partial<CanvasElement>) => void;
    removeElement: (id: string) => void;
    setSelectedId: (id: string | null) => void;
    setActiveTool: (tool: EditorState["activeTool"]) => void;
}

export const useEditorStore = create<EditorState>((set) => ({
    elements: [],
    selectedId: null,
    activeTool: null,

    setElements: (elements) => set({ elements }),
    addElement: (element) => set((state) => ({
        elements: [...state.elements, element],
        selectedId: element.id
    })),
    updateElement: (id, newProps) => set((state) => ({
        elements: state.elements.map((el) => el.id === id ? { ...el, ...newProps } : el)
    })),
    removeElement: (id) => set((state) => ({
        elements: state.elements.filter((el) => el.id !== id),
        selectedId: state.selectedId === id ? null : state.selectedId,
    })),
    setSelectedId: (id) => set({ selectedId: id }),
    setActiveTool: (tool) => set({ activeTool: tool }),
}));