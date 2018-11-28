<?php

/*
* Working with API
*/

namespace app\models;

use Yii;
use yii\base\Model;

/**
 * LoginForm is the model behind the login form.
 */
class Library extends Model
{
    private static function getCurl($url){
        $ch = curl_init();
        curl_setopt($ch, CURLOPT_URL, $url);
        curl_setopt($ch, CURLOPT_RETURNTRANSFER, 1);
        $result = curl_exec($ch);
        curl_close($ch);
        return $result;
    }

    public static function getBooks($name = null)
    {
        if (!empty($name)) {
            $result = self::getCurl('http://localhost:5051/api/books/'.$name);
            return $result;
        } else {
            $result = self::getCurl('http://localhost:5051/api/books/');
	        return $result;
        }
    }
}
