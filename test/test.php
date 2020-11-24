<?php



$url = 'http://0.0.0.0:5000/processing-image-by-url';  //调用接口的平台服务地址

$post_data = "imageFileUrl=https://pic4.zhimg.com/50/v2-76a08e1141cb62451ea89d32ab84ebab_hd.jpg&metaId=12345678" . time();

//post数据

function postProcessing($url,$post_data){

    $curl = curl_init();

    curl_setopt($curl, CURLOPT_URL, $url);

    curl_setopt($curl, CURLOPT_POST, 1 );

    curl_setopt($curl, CURLOPT_POSTFIELDS, $post_data);

    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);

    curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, FALSE);

    curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, FALSE);

    $response = curl_exec($curl);

    var_dump($response);

    $result = json_decode($response,true);

    $error = curl_error($curl);

    return $error ? $error : $result;

}

$result = postProcessing($url,$post_data);

var_dump($result);




$url = 'http://0.0.0.0:5000/query-image-by-url';  //调用接口的平台服务地址

$post_data = "imageFileUrl=https://pic4.zhimg.com/50/v2-76a08e1141cb62451ea89d32ab84ebab_hd.jpg";

//post数据

function postQuery($url,$post_data){

    $curl = curl_init();

    curl_setopt($curl, CURLOPT_URL, $url);

    curl_setopt($curl, CURLOPT_POST, 1 );

    curl_setopt($curl, CURLOPT_POSTFIELDS, $post_data);

    curl_setopt($curl, CURLOPT_RETURNTRANSFER, 1);

    curl_setopt($curl, CURLOPT_SSL_VERIFYPEER, FALSE);

    curl_setopt($curl, CURLOPT_SSL_VERIFYHOST, FALSE);

    $response = curl_exec($curl);

    $result = json_decode($response,true);

    $error = curl_error($curl);

    return $error ? $error : $result;

}

$result = postQuery($url,$post_data);

var_dump($result);




