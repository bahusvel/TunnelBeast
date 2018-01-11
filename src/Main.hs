{-# LANGUAGE OverloadedStrings #-}
module Main where

import           Control.Applicative        ((<$>))
import           Control.Monad
import qualified Data.ByteString.Char8      as BS
import qualified Data.ByteString.Lazy.Char8 as LBS
import           Data.IP
import           Data.Maybe
import           Data.Monoid
import           Network.HTTP.Types         (status200)
import           Network.Socket
import           Network.Wai
import           Network.Wai.Handler.Warp

import qualified Data.Text.Format           as TF
import           System.Exit
import           System.IO
import           System.Process

data Route = Route {
    srcip   :: String,
    dstip   :: String,
    srcport :: String,
    dstport :: String
} deriving Show

addRouteCommand = "iptables -t nat -A PREROUTING -i {} -p {} -s {} -j DNAT --to-destination {}:{} --dport {}"

main = do
    let port = 3000
    putStrLn $ "Listening on port " ++ show port
    run port app

app req respond = case rawPathInfo req of
    "/"    -> index respond
    "/add" -> addHandler req respond

addRoute :: Route -> IO ExitCode
addRoute r = do
    tcp <- spawnCommand $ show $ TF.format addRouteCommand ("eth0" :: String , "tcp" :: String, srcip r, dstip r, dstport r, srcport r)
    tcpResult <- waitForProcess tcp
    case tcpResult of
        ExitFailure code -> return $ ExitFailure code
    udp <- spawnCommand $ show $ TF.format addRouteCommand ("eth0" :: String , "udp" :: String, srcip r, dstip r, dstport r, srcport r)
    udpResult <- waitForProcess udp
    case udpResult of
        ExitFailure code -> return $ ExitFailure code
        _                -> return ExitSuccess


addHandler req respond =
    let q = queryString req
        SockAddrInet _ raddr = remoteHost req
        uname = join $ lookup "username" q
        passwd = join $ lookup "password" q
        srcip = join $ lookup "sourceip" q
        dstip = join $ lookup "internalip" q
        dstport = join $ lookup "internalport" q
        srcport = join $ lookup "externalport" q
    in case (uname, passwd, dstip, dstport, srcport) of
        (Just u, Just p, Just di, Just dp, Just sp) -> let route = Route { srcip = fromMaybe (show $ fromHostAddress raddr) (BS.unpack <$> srcip), dstip = BS.unpack di, dstport = BS.unpack dp, srcport = BS.unpack sp}
            in do
                addRoute route
                respond $ responseLBS status200 [ ("Content-Type", "text/plain") ] $ LBS.pack $ "Query parameter: " ++ (show (u, p, route))
        _ -> respond $ responseLBS status200 [ ("Content-Type", "text/plain") ] "ERROR INPUT"



yay respond = do
    responseLBS
        status200
        [ ("Content-Type", "text/plain") ]
        "yay\n"


index respond = respond $ responseFile
    status200
    [ ("Content-Type", "text/html") ]
    "html/index.html"
    Nothing
