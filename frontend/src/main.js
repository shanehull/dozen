import { CalcService } from "../bindings/github.com/shanehull/dozen";

// NOP is any op string the engine does not recognise: under an f/g prefix it
// falls through to the default case, which simply clears the prefix. We use it
// for shifted labels that have no engine implementation.
const NOP = "NOP";

// Each key: r/c grid position, prim (face label), g (blue g-label), f (gold
// f-label). op = unprefixed engine op. fop/gop = engine op to send while the
// f/g prefix is armed (the engine routes it through prefixedF / prefixedG).
// fixArg: digit value used for f+digit (FIX n).
const KEYS = [
  // row 1
  { r:1,c:1,  prim:"n",   f:"AMORT", g:"12ˣ", op:"n",   fop:"AMORT", gop:NOP,  tvm:true },
  { r:1,c:2,  prim:"i",   f:"INT",   g:"12÷", op:"i",   fop:"INT",   gop:NOP,  tvm:true },
  { r:1,c:3,  prim:"PV",  f:"NPV",   g:"CFo", op:"PV",  fop:"NPV",   gop:"CF0",  tvm:true },
  { r:1,c:4,  prim:"PMT", f:"RND",   g:"CFj", op:"PMT", fop:NOP,     gop:"CFj",  tvm:true },
  { r:1,c:5,  prim:"FV",  f:"IRR",   g:"Nⱼ",  op:"FV",  fop:"IRR",   gop:"Nj",   tvm:true },
  { r:1,c:6,  prim:"CHS", f:"",      g:"DATE",op:"CHS", fop:NOP,     gop:"DATE" },
  { r:1,c:7,  prim:"7",   f:"",      g:"BEG", op:"7",   fop:"FIX",   gop:"BEG", fixArg:7 },
  { r:1,c:8,  prim:"8",   f:"",      g:"END", op:"8",   fop:"FIX",   gop:"END", fixArg:8 },
  { r:1,c:9,  prim:"9",   f:"",      g:"MEM", op:"9",   fop:"FIX",   gop:"MEM", fixArg:9 },
  { r:1,c:10, prim:"÷",   f:"",     g:"↩",  op:"÷",   fop:NOP,     gop:"÷" },


  // row 2
  { r:2,c:1,  prim:"yˣ",  f:"PRICE", g:"√x",  op:"yˣ",  fop:"PRICE", gop:"√x", primClass:"sm" },
  { r:2,c:2,  prim:"1/x", f:"YTM",   g:"eˣ",  op:"1/x", fop:"YTM",   gop:"eˣ", primClass:"sm" },
  { r:2,c:3,  prim:"%T",  f:"SL",    g:"LN",  op:"%T",  fop:"SL",    gop:"LN", primClass:"sm" },
  { r:2,c:4,  prim:"Δ%",  f:"SOYD",  g:"FRAC",op:"%CHG",fop:"SOYD",  gop:"FRAC", primClass:"sm" },
  { r:2,c:5,  prim:"%",   f:"DB",    g:"INTG",op:"%",   fop:"DB",    gop:"INTG" },
  { r:2,c:6,  prim:"EEX", f:"",      g:"ΔDYS",op:"EEX", fop:NOP,     gop:"DYS", primClass:"sm" },
  { r:2,c:7,  prim:"4",   f:"",      g:"D.MY",op:"4",   fop:"FIX",   gop:"D.MY", fixArg:4 },
  { r:2,c:8,  prim:"5",   f:"",      g:"M.DY",op:"5",   fop:"FIX",   gop:"M.DY", fixArg:5 },
  { r:2,c:9,  prim:"6",   f:"",      g:"x̄w", op:"6",   fop:"FIX",   gop:"x̄w", fixArg:6 },
  { r:2,c:10, prim:"×",   f:"",     g:"x²",  op:"×",   fop:NOP,     gop:"×" },

  // row 3
  { r:3,c:1,  prim:"R/S", f:"P/R",   g:"PSE", op:"R/S", fop:"P/R",   gop:"PSE", primClass:"sm" },
  { r:3,c:2,  prim:"SST", f:"Σ",     g:"BST", op:"SST", fop:"CLEAR Σ", gop:"BST", primClass:"sm" },
  { r:3,c:3,  prim:"R↓",  f:"PRGM",  g:"GTO", op:"R↓",  fop:"CLEAR PRGM", gop:"GTO" },
  { r:3,c:4,  prim:"x↔y", f:"FIN",   g:"x≤y", op:"x↔y", fop:"CLEAR FIN", gop:"x≤y", primClass:"xs" },
  { r:3,c:5,  prim:"CLx", f:"REG",   g:"x=0", op:"CLx", fop:"CLEAR REG", gop:"x=0", primClass:"sm" },
  { enter:true, r:3, c:6, prim:"ENTER", f:"PREFIX", g:"",   op:"ENTER", fop:"CLEAR PREFIX", gop:NOP },
  { r:3,c:7,  prim:"1",   f:"",      g:"x̂,r", op:"1",  fop:NOP,     gop:NOP },
  { r:3,c:8,  prim:"2",   f:"",      g:"ŷ,r", op:"2",  fop:NOP,     gop:"ŷ,r" },
  { r:3,c:9,  prim:"3",   f:"",      g:"n!",  op:"3",   fop:NOP,     gop:"n!" },
  { r:3,c:10, prim:"−",   f:"",     g:"←",   op:"−",   fop:NOP,     gop:"−" },
  // row 4
  { r:4,c:1,  prim:"ON",  f:"",      g:"",    op:"ON",  fop:NOP,     gop:NOP, primClass:"sm" },
  { r:4,c:2,  prim:"f",   f:"",      g:"",    op:"f",   special:"f" },
  { r:4,c:3,  prim:"g",   f:"",      g:"",    op:"g",   special:"g" },
  { r:4,c:4,  prim:"STO", f:"",      g:"",    op:"STO", special:"STO", primClass:"sm" },
  { r:4,c:5,  prim:"RCL", f:"",      g:"",    op:"RCL", special:"RCL", primClass:"sm" },
  { r:4,c:7,  prim:"0",   f:"",      g:"x̄",  op:"0",   fop:"FIX",   gop:"x̄", fixArg:0 },
  { r:4,c:8,  prim:".",   f:"",      g:"s",   op:".",   fop:NOP,     gop:"s" },
  { r:4,c:9,  prim:"Σ+",  f:"",      g:"Σ−",  op:"Σ+",  fop:NOP,     gop:NOP, primClass:"sm" },
  { r:4,c:10, prim:"+",   f:"",      g:"LSTx", op:"+",   fop:NOP,     gop:"+" },
];

// Groups drawn as gold brackets above the keys.
const GROUPS = [
  { row:2, from:1, to:2, cap:"BOND" },
  { row:2, from:3, to:5, cap:"DEPRECIATION" },
  { row:3, from:2, to:6, cap:"CLEAR" },
];

const gridEl    = document.getElementById("grid");
const bracketEl = document.getElementById("brackets");
const digitsEl  = document.getElementById("digits");
const signEl    = document.getElementById("sign");
const annunEl   = document.getElementById("annun");

let armed = null;          // null | "f" | "g" | "STO" | "RCL"
let entering = false;      // is a number currently being keyed in?
const byOp = {};           // op -> button element (for keyboard highlight)

// ---- backend calls -------------------------------------------------------
function callStep(op, arg) {
  return CalcService.PressKey({ op, arg: arg || 0, argS: op }).then(state => {
    CalcService.Save().then(j => localStorage.setItem("dozen-state", j));
    return state;
  });
}

async function sequence(steps) {
  let state;
  for (const s of steps) state = await callStep(s.op, s.arg);
  if (state) applyState(state);
}

// ---- display -------------------------------------------------------------
const ANNUN = ["f", "g", "BEGIN", "RAD", "D.MY", "PRGM"];

function applyState(s) {
  const d = (s && s.display) || {};
  const sign = d.sign && d.sign.trim() ? d.sign : "";
  signEl.textContent = sign;
  digitsEl.textContent = d.mantissa && d.mantissa.trim() ? d.mantissa.trim() : "0.00";
  const flags = (d.flags || []).map(f => String(f).toUpperCase());
  annunEl.innerHTML = ANNUN.map(a => {
    const on = flags.includes(a.toUpperCase());
    return `<span class="${on ? "on" : ""}">${a}</span>`;
  }).join("");
}

// ---- input model ---------------------------------------------------------
const DIGITS = new Set(["0","1","2","3","4","5","6","7","8","9"]);

function isEntryKey(op) {
  return DIGITS.has(op) || op === "." || op === "EEX" || op === "CHS";
}

function setArmed(next) {
  armed = next;
  refreshArmed();
}

function refreshArmed() {
  for (const el of gridEl.querySelectorAll(".key")) el.classList.remove("armed");
  if (armed && byOp[armed]) byOp[armed].classList.add("armed");
}

function handle(def) {
  const op = def.op;

  // f / g / STO / RCL are prefix keys handled by the engine.
  if (def.special) {
    callStep(op).then(applyState);
    setArmed(armed === def.special ? null : def.special);
    entering = false;
    return;
  }

  if (armed === "f") {
    const fop = def.fop || NOP;
    if (fop === "FIX") sequence([{ op:"FIX", arg: def.fixArg }]);
    else sequence([{ op: fop }]);
    setArmed(null);
    entering = false;
    return;
  }

  if (armed === "g") {
    sequence([{ op: def.gop || NOP }]);
    setArmed(null);
    entering = false;
    return;
  }

  if (armed === "STO" || armed === "RCL") {
    // register index comes from a digit key
    const arg = DIGITS.has(op) ? parseInt(op, 10) : 0;
    sequence([{ op, arg }]);
    setArmed(null);
    entering = false;
    return;
  }

  // No prefix armed.
  // TVM keys: the engine's StackLift handles store vs solve.
  if (def.tvm) {
    sequence([{ op }]);
    entering = false;
    return;
  }

  // ON acts as a soft clear.
  if (op === "ON") { sequence([{ op: "CLx" }]); entering = false; return; }

  const arg = DIGITS.has(op) ? parseInt(op, 10) : 0;
  sequence([{ op, arg }]);
  entering = op === "ENTER" ? false : isEntryKey(op);
}

// ---- build DOM -----------------------------------------------------------
function buildKeys() {
  gridEl.innerHTML = "";
  for (const def of KEYS) {
    const btn = document.createElement("button");
    btn.className = "key";
    if (def.enter) btn.classList.add("enter");
    if (def.special === "f") btn.classList.add("fkey");
    if (def.special === "g") btn.classList.add("gkey");
    if (!def.enter) {
      btn.style.gridColumn = String(def.c);
      btn.style.gridRow = String(def.r);
    }

    if (def.f) {
      const fl = document.createElement("span");
      fl.className = "flabel";
      fl.textContent = def.f;
      btn.appendChild(fl);
    }
    const prim = document.createElement("span");
    prim.className = "prim" + (def.primClass ? " " + def.primClass : "");
    prim.textContent = def.prim;
    btn.appendChild(prim);

    if (def.g) {
      const gl = document.createElement("span");
      gl.className = "glabel";
      gl.textContent = def.g;
      btn.appendChild(gl);
    }

    btn.addEventListener("click", () => handle(def));
    byOp[def.op] = btn;
    gridEl.appendChild(btn);
  }
  gridEl.appendChild(bracketEl); // re-attach overlay (innerHTML wipe removed it)
}

function drawBrackets() {
  bracketEl.innerHTML = "";
  // Use offsetLeft/offsetTop (layout coordinates, unaffected by the fit-scale
  // CSS transform) relative to #grid, which is the overlay's offset parent.
  const cellEl = (r, c) => {
    const def = KEYS.find(k => k.r === r && k.c === c);
    return def ? byOp[def.op] : null;
  };
  for (const g of GROUPS) {
    const a = cellEl(g.row, g.from);
    const b = cellEl(g.row, g.to);
    if (!a || !b) continue;
    const div = document.createElement("div");
    div.className = "bracket";
    const left = a.offsetLeft;
    const right = b.offsetLeft + b.offsetWidth;
    div.style.left = left + "px";
    div.style.width = (right - left) + "px";
    div.style.top = (a.offsetTop - 23) + "px"; // just above the gold f-labels
    const segL = document.createElement("span");
    segL.className = "seg l";
    const cap = document.createElement("span");
    cap.className = "cap";
    cap.textContent = g.cap;
    const segR = document.createElement("span");
    segR.className = "seg r";
    div.append(segL, cap, segR);
    bracketEl.appendChild(div);
  }
}

// ---- scaling to fit window ----------------------------------------------
function fitScale() {
  const calc = document.getElementById("calc");
  const scaler = document.getElementById("scaler");
  const pad = 16;
  const sx = (window.innerWidth - pad) / calc.offsetWidth;
  const sy = (window.innerHeight - pad) / calc.offsetHeight;
  const s = Math.min(sx, sy);
  scaler.style.transform = `scale(${s})`;
}

// ---- keyboard ------------------------------------------------------------
const KEYMAP = {
  "+": "+", "-": "−", "*": "×", "/": "÷",
  "Enter": "ENTER", ".": ".",
  "n": "n", "i": "i",
  "p": "PV", "m": "PMT", "v": "FV",
  "s": "STO", "r": "RCL",
  "x": "x↔y",
};

document.addEventListener("keydown", e => {
  const k = e.key;
  let op = null;
  if (DIGITS.has(k)) op = k;
  else if (KEYMAP[k]) op = KEYMAP[k];
  else if (k === "Escape" || k === "Backspace") op = "CLx";
  else if (k === "%") op = "%";
  else if (k === "f" || k === "F") op = "f";
  else if (k === "g" || k === "G") op = "g";
  if (!op) return;
  e.preventDefault();
  const def = KEYS.find(d => d.op === op);
  if (!def) return;
  handle(def);
  const btn = byOp[op];
  if (btn) { btn.classList.add("pressed"); setTimeout(() => btn.classList.remove("pressed"), 90); }
});

// ---- init ----------------------------------------------------------------
buildKeys();
requestAnimationFrame(() => { drawBrackets(); fitScale(); });
window.addEventListener("resize", () => { drawBrackets(); fitScale(); });

const saved = localStorage.getItem("dozen-state");
if (saved) {
  CalcService.Load(saved).then(() => CalcService.GetState().then(applyState));
} else {
  CalcService.GetState().then(applyState).catch(() =>
    applyState({ display: { mantissa: "0.00", sign: "", flags: [] } }));
}
