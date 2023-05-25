import {
  Component,
  createEffect,
  createSignal,
  JSX,
  onCleanup,
} from "solid-js";

const App: Component = () => {
  const [eventSource, setEventSource] = createSignal<EventSource | undefined>();
  const [time, setTime] = createSignal("");
  const [id, setId] = createSignal("");

  async function handleGetTime() {
    const res = await fetch(`http://localhost:3500/time/${id()}`);
    if (res.status !== 200) {
      console.log("Could not connect to the server");
    } else {
      console.log("OK");
    }
  }

  function handleChange(e: InputEvent) {
    setId((e.currentTarget as HTMLInputElement).value);
  }

  function handleConnect() {
    const ev = new EventSource(`http://localhost:3500/sse/${id()}`);
    ev.addEventListener("timeEvent", (e) => {
      console.log({ event: e.type });
      console.log({ data: e.data });

      setTime(e.data);
    });

    setEventSource(ev);
  }

  onCleanup(() => {
    eventSource()?.close();
  });

  return (
    <main>
      Time: {time()}
      <button onClick={handleGetTime}>Get time</button>
      <input type="text" onInput={handleChange} value={id()} />
      <button onClick={handleConnect}>Connect</button>
    </main>
  );
};

export default App;
