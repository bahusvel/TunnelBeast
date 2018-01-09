{-# LANGUAGE OverloadedStrings #-}
module Main where

import qualified Data.ByteString.Char8    as BS
import           Data.Monoid
import           Network.HTTP.Types       (status200)
import           Network.Wai
import           Network.Wai.Handler.Warp

main = do
    let port = 3000
    putStrLn $ "Listening on port " ++ show port
    run port app

app req respond = respond $ case rawPathInfo req of
    "/" -> index

yay respond = do
    responseLBS
        status200
        [ ("Content-Type", "text/plain") ]
        "yay\n"


index :: Response
index = responseFile
    status200
    [ ("Content-Type", "text/html") ]
    "html/index.html"
    Nothing
