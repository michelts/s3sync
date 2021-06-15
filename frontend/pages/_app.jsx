import '../styles/globals.css';

function MyApp({ Component, pageProps }) {
  return (
    <>
      <div className="container w-full mx-auto p-3">
        <Component {...pageProps} />
      </div>
    </>
  )
}

export default MyApp
