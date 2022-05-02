import {BrowserRouter, Routes, Route} from 'react-router-dom';
import Home from 'pages/Home';
import Test from 'pages/Test';

const Router = (): JSX.Element => {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home path="home" />} />
        <Route path="/test/:id" element={<Test path="test-details" />} />
      </Routes>
    </BrowserRouter>
  );
};

export default Router;