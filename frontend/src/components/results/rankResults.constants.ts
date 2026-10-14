export const GRID_COLS =
  "grid-cols-[60px_minmax(0,1.5fr)_minmax(0,1fr)] " +
  "@[456px]:grid-cols-[60px_minmax(0,1.5fr)_minmax(0,1fr)_minmax(0,1fr)] " +
  "@[736px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(2,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[996px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(3,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1256px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(4,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1516px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(5,minmax(0,1fr))_minmax(0,1fr)] " +
  "@[1776px]:grid-cols-[60px_minmax(0,1.5fr)_repeat(6,minmax(0,1fr))_minmax(0,1fr)]";

export const ROUND_VISIBILITY = [
  "hidden @[1776px]:flex",
  "hidden @[1516px]:flex",
  "hidden @[1256px]:flex",
  "hidden @[996px]:flex",
  "hidden @[736px]:flex",
  "hidden @[456px]:flex",
];