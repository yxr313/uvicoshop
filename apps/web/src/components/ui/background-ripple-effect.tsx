"use client";
import React, { useMemo, useState } from "react";
import { cn } from "@/lib/utils";

export const BackgroundRippleEffect = ({
  rows = 25,
  cols = 25,
  cellSize = 60,
}: {
  rows?: number;
  cols?: number;
  cellSize?: number;
}) => {
  const [clickedCell, setClickedCell] = useState<{
    row: number;
    col: number;
  } | null>(null);
  const [rippleKey, setRippleKey] = useState(0);

  return (
    <div className="absolute inset-0 h-full w-full overflow-hidden bg-[#0a0a0a] flex items-center justify-center">
      <DivGrid
        key={`base-${rippleKey}`}
        rows={rows}
        cols={cols}
        cellSize={cellSize}
        clickedCell={clickedCell}
        onCellClick={(row, col) => {
          setClickedCell({ row, col });
          setRippleKey((k) => k + 1);
        }}
        interactive
      />
    </div>
  );
};

type DivGridProps = {
  rows: number;
  cols: number;
  cellSize: number;
  clickedCell: { row: number; col: number } | null;
  onCellClick?: (row: number, col: number) => void;
  interactive?: boolean;
};

type CellStyle = React.CSSProperties & {
  ["--delay"]?: string;
  ["--duration"]?: string;
  ["--animate-cell-ripple"]?: string;
};

const DivGrid = ({
  rows,
  cols,
  cellSize,
  clickedCell,
  onCellClick = () => {},
  interactive = true,
}: DivGridProps) => {
  const cells = useMemo(
    () => Array.from({ length: rows * cols }, (_, idx) => idx),
    [rows, cols],
  );

  const gridStyle: React.CSSProperties = {
    display: "grid",
    gridTemplateColumns: `repeat(${cols}, ${cellSize}px)`,
    gridTemplateRows: `repeat(${rows}, ${cellSize}px)`,
  };

  return (
    <div className="relative z-[3]" style={gridStyle}>
      {cells.map((idx) => {
        const rowIdx = Math.floor(idx / cols);
        const colIdx = idx % cols;
        const distance = clickedCell
          ? Math.hypot(clickedCell.row - rowIdx, clickedCell.col - colIdx)
          : 0;
        const delay = clickedCell ? Math.max(0, distance * 50) : 0; // ms
        const duration = 400 + distance * 80; // ms

        const style: CellStyle = clickedCell
          ? {
              "--delay": `${delay}ms`,
              "--duration": `${duration}ms`,
              // Inyectamos la animación directamente si hubo un clic
              animation: `cell-ripple var(--duration) ease-out forwards var(--delay)`,
            }
          : {};

        return (
          <div
            key={idx}
            className={cn(
              "cell relative border-[0.5px] transition-colors duration-300 will-change-transform",
              !clickedCell && "border-white/5 hover:bg-[rgba(212,215,0,0.08)] hover:border-[#d4d700]/30",
              !interactive && "pointer-events-none"
            )}
            style={style}
            onClick={
              interactive ? () => onCellClick(rowIdx, colIdx) : undefined
            }
          />
        );
      })}
    </div>
  );
};
